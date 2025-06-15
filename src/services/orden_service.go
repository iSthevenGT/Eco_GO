package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

type OrdenService struct {
	carritoService *CarritoService
	fechaService   *FechaService
	pagoService    *PagoService
}

func NewOrdenService() *OrdenService {
	return &OrdenService{
		carritoService: NewCarritoService(),
		fechaService:   NewFechaService(),
		pagoService:    NewPagoService(),
	}
}

func (s *OrdenService) Crear(consumidorID uint, ordenData map[string]interface{}) (*models.Orden, error) {
	// Obtener carrito del consumidor
	carrito := s.carritoService.ObtenerCarrito(consumidorID)
	if len(carrito.Productos) == 0 {
		return nil, errors.New("el carrito está vacío")
	}

	// Verificar stock de los productos
	for _, productoDTO := range carrito.Productos {
		var producto models.Producto
		if err := database.DB.First(&producto, productoDTO.ID).Error; err != nil {
			return nil, errors.New("producto no encontrado")
		}
		if producto.Stock < productoDTO.Cantidad {
			return nil, errors.New("stock insuficiente para el producto: " + producto.Nombre)
		}
	}

	// Obtener el consumidor
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, consumidorID).Error; err != nil {
		return nil, errors.New("consumidor no encontrado")
	}

	// Iniciar transacción
	tx := database.DB.Begin()

	// Crear fecha de la orden
	fechaData := ordenData["fechaOrden"].(map[string]interface{})
	fecha := models.Fecha{
		Anio: int(fechaData["anio"].(float64)),
		Mes:  int(fechaData["mes"].(float64)),
		Dia:  int(fechaData["dia"].(float64)),
	}
	if err := tx.Create(&fecha).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Crear orden
	orden := models.Orden{
		ConsumidorID: consumidorID,
		FechaOrdenID: fecha.ID,
		MontoTotal:   carrito.Total,
		EstadoOrden:  models.EstadoPendiente,
	}

	// Asignar dirección de entrega si existe
	if direccionData, ok := ordenData["direccionEntrega"].(map[string]interface{}); ok {
		if idDireccion, exists := direccionData["idDireccion"]; exists && idDireccion != nil {
			direccionID := uint(idDireccion.(float64))
			orden.DireccionEntregaID = &direccionID
		}
	}

	if err := tx.Create(&orden).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Crear pago
	pagoData := ordenData["pago"].(map[string]interface{})
	pago := models.Pago{
		OrdenID:     orden.ID,
		MetodoPago:  models.MetodoPago(pagoData["metodoPago"].(string)),
		EstadoPago:  models.EstadoPago(pagoData["estadoPago"].(string)),
		MontoPagado: carrito.Total,
		FechaPagoID: fecha.ID,
	}
	if err := tx.Create(&pago).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Crear orden-productos y actualizar stock
	for _, productoDTO := range carrito.Productos {
		var producto models.Producto
		if err := tx.First(&producto, productoDTO.ID).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Actualizar stock
		producto.Stock -= productoDTO.Cantidad
		if err := tx.Save(&producto).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Crear asociación orden-producto
		ordenProducto := models.OrdenProducto{
			OrdenID:    orden.ID,
			ProductoID: producto.ID,
			Cantidad:   productoDTO.Cantidad,
		}
		if err := tx.Create(&ordenProducto).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Limpiar carrito
	s.carritoService.LimpiarCarrito(consumidorID)

	// Cargar orden completa
	if err := database.DB.Preload("Consumidor").Preload("FechaOrden").Preload("DireccionEntrega").Preload("Productos.Producto").Preload("Pago").First(&orden, orden.ID).Error; err != nil {
		return nil, err
	}

	// Auto-confirmar si el pago está aprobado
	if pago.EstadoPago == models.PagoAprobado {
		go func() {
			time.Sleep(1 * time.Second)
			s.Confirmar(orden.ID)
		}()
	}

	return &orden, nil
}

func (s *OrdenService) ObtenerPorID(consumidorID, ordenID uint) (*models.Orden, error) {
	var orden models.Orden
	if err := database.DB.Where("id = ? AND consumidor_id = ?", ordenID, consumidorID).
		Preload("Consumidor").Preload("FechaOrden").Preload("DireccionEntrega").
		Preload("Productos.Producto").Preload("Pago").First(&orden).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("orden no encontrada")
		}
		return nil, err
	}
	return &orden, nil
}

func (s *OrdenService) ObtenerTodosPorComerciante(comercianteID uint) ([]models.Orden, error) {
	var ordenes []models.Orden
	if err := database.DB.Joins("JOIN orden_productos ON ordenes.id = orden_productos.orden_id").
		Joins("JOIN productos ON orden_productos.producto_id = productos.id").
		Where("productos.comerciante_id = ?", comercianteID).
		Preload("Consumidor").Preload("FechaOrden").
		Preload("Productos", "productos.comerciante_id = ?", comercianteID).
		Preload("Productos.Producto").
		Distinct().Find(&ordenes).Error; err != nil {
		return nil, err
	}
	return ordenes, nil
}

func (s *OrdenService) ProductosPorIDAndComerciante(comercianteID, ordenID uint) ([]models.OrdenProducto, error) {
	var productos []models.OrdenProducto
	if err := database.DB.Joins("JOIN productos ON orden_productos.producto_id = productos.id").
		Where("orden_productos.orden_id = ? AND productos.comerciante_id = ?", ordenID, comercianteID).
		Preload("Producto").Find(&productos).Error; err != nil {
		return nil, err
	}
	return productos, nil
}

func (s *OrdenService) Confirmar(ordenID uint) error {
	var orden models.Orden
	if err := database.DB.First(&orden, ordenID).Error; err != nil {
		return errors.New("orden no encontrada")
	}

	if orden.EstadoOrden != models.EstadoPendiente {
		return errors.New("la orden no está en estado pendiente")
	}

	orden.EstadoOrden = models.EstadoConfirmada
	return database.DB.Save(&orden).Error
}

func (s *OrdenService) Cancelar(consumidorID, ordenID uint) error {
	var orden models.Orden
	if err := database.DB.Where("id = ? AND consumidor_id = ?", ordenID, consumidorID).First(&orden).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("orden no encontrada")
		}
		return err
	}

	if orden.EstadoOrden == models.EstadoCancelada {
		return errors.New("la orden ya está cancelada")
	}

	if orden.EstadoOrden == models.EstadoReembolsada {
		return errors.New("no se puede cancelar una orden reembolsada")
	}

	orden.EstadoOrden = models.EstadoCancelada
	return database.DB.Save(&orden).Error
}
