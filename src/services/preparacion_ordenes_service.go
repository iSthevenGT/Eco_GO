package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/dto"
	"Eco_GO/src/models"
	"errors"
	"sync"
)

type PreparacionOrdenesService struct {
	ordenesPrep  map[uint]*dto.OrdenesPrepDTO
	mutex        sync.RWMutex
	ordenService *OrdenService
}

func NewPreparacionOrdenesService() *PreparacionOrdenesService {
	return &PreparacionOrdenesService{
		ordenesPrep:  make(map[uint]*dto.OrdenesPrepDTO),
		ordenService: NewOrdenService(),
	}
}

func (s *PreparacionOrdenesService) ObtenerOrdenesPrep(comercianteID uint) *dto.OrdenesPrepDTO {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if ordenes, exists := s.ordenesPrep[comercianteID]; exists {
		return ordenes
	}

	newOrdenes := &dto.OrdenesPrepDTO{
		Ordenes: []interface{}{},
	}
	s.ordenesPrep[comercianteID] = newOrdenes
	return newOrdenes
}

func (s *PreparacionOrdenesService) AgregarOrden(comercianteID, ordenID uint) (*models.Orden, error) {
	// Confirmar la orden
	if err := s.ordenService.Confirmar(ordenID); err != nil {
		return nil, err
	}

	// Obtener la orden confirmada
	var orden models.Orden
	if err := database.DB.Preload("Consumidor").Preload("FechaOrden").
		Preload("Productos.Producto").Preload("Pago").First(&orden, ordenID).Error; err != nil {
		return nil, errors.New("orden no encontrada")
	}

	// Agregar a la lista de preparación
	ordenes := s.ObtenerOrdenesPrep(comercianteID)
	ordenes.Ordenes = append(ordenes.Ordenes, orden)

	return &orden, nil
}

func (s *PreparacionOrdenesService) ObtenerOrdenes(comercianteID uint) *dto.OrdenesPrepDTO {
	return s.ObtenerOrdenesPrep(comercianteID)
}

func (s *PreparacionOrdenesService) ObtenerOrden(comercianteID, ordenID uint) (interface{}, error) {
	ordenes := s.ObtenerOrdenesPrep(comercianteID)

	for _, orden := range ordenes.Ordenes {
		if ordenModel, ok := orden.(models.Orden); ok {
			if ordenModel.ID == ordenID {
				return ordenModel, nil
			}
		}
	}

	return nil, errors.New("orden no encontrada en preparación")
}
