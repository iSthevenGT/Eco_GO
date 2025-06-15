package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/dto"
	"Eco_GO/src/models"
	"errors"
	"sync"
)

type CarritoService struct {
	carritos map[uint]*dto.CarritoDTO
	mutex    sync.RWMutex
}

func NewCarritoService() *CarritoService {
	return &CarritoService{
		carritos: make(map[uint]*dto.CarritoDTO),
	}
}

func (s *CarritoService) ObtenerCarrito(consumidorID uint) *dto.CarritoDTO {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if carrito, exists := s.carritos[consumidorID]; exists {
		return carrito
	}

	newCarrito := &dto.CarritoDTO{
		Productos: []dto.ProductoDTO{},
		Total:     0,
	}
	s.carritos[consumidorID] = newCarrito
	return newCarrito
}

func (s *CarritoService) AgregarProducto(consumidorID, productoID uint, cantidad int) error {
	if cantidad <= 0 {
		return errors.New("cantidad no válida")
	}

	// Verificar que el consumidor existe
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, consumidorID).Error; err != nil {
		return errors.New("consumidor no encontrado")
	}

	// Verificar que el producto existe
	var producto models.Producto
	if err := database.DB.First(&producto, productoID).Error; err != nil {
		return errors.New("producto no encontrado")
	}

	carrito := s.ObtenerCarrito(consumidorID)

	// Buscar si el producto ya está en el carrito
	for i := range carrito.Productos {
		if carrito.Productos[i].ID == productoID {
			carrito.Productos[i].Cantidad += cantidad
			s.calcularTotal(carrito)
			return nil
		}
	}

	// Agregar nuevo producto
	nuevoProducto := dto.ProductoDTO{
		ID:       producto.ID,
		Imagen:   producto.Imagen,
		Nombre:   producto.Nombre,
		Precio:   producto.Precio,
		Cantidad: cantidad,
	}

	carrito.Productos = append(carrito.Productos, nuevoProducto)
	s.calcularTotal(carrito)
	return nil
}

func (s *CarritoService) EliminarProducto(consumidorID, productoID uint) error {
	carrito := s.ObtenerCarrito(consumidorID)

	for i, producto := range carrito.Productos {
		if producto.ID == productoID {
			carrito.Productos = append(carrito.Productos[:i], carrito.Productos[i+1:]...)
			s.calcularTotal(carrito)
			return nil
		}
	}

	return errors.New("producto no encontrado en el carrito")
}

func (s *CarritoService) CambiarCantidad(consumidorID, productoID uint, cantidad int) error {
	if cantidad <= 0 {
		return errors.New("cantidad no válida")
	}

	carrito := s.ObtenerCarrito(consumidorID)

	for i := range carrito.Productos {
		if carrito.Productos[i].ID == productoID {
			carrito.Productos[i].Cantidad = cantidad
			s.calcularTotal(carrito)
			return nil
		}
	}

	return errors.New("producto no encontrado en el carrito")
}

func (s *CarritoService) LimpiarCarrito(consumidorID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.carritos, consumidorID)
}

func (s *CarritoService) CalcularTotal(consumidorID uint) int {
	carrito := s.ObtenerCarrito(consumidorID)
	return carrito.Total
}

func (s *CarritoService) calcularTotal(carrito *dto.CarritoDTO) {
	total := 0
	for _, producto := range carrito.Productos {
		total += producto.Precio * producto.Cantidad
	}
	carrito.Total = total
}
