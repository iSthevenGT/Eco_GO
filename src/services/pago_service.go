package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"

	"gorm.io/gorm"
)

type PagoService struct{}

func NewPagoService() *PagoService {
	return &PagoService{}
}

func (s *PagoService) Crear(pago models.Pago) (*models.Pago, error) {
	if err := database.DB.Create(&pago).Error; err != nil {
		return nil, err
	}
	return &pago, nil
}

func (s *PagoService) ObtenerPorID(id uint) (*models.Pago, error) {
	var pago models.Pago
	if err := database.DB.First(&pago, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pago no encontrado")
		}
		return nil, err
	}
	return &pago, nil
}

func (s *PagoService) Actualizar(pago models.Pago) (*models.Pago, error) {
	var existente models.Pago
	if err := database.DB.First(&existente, pago.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pago no encontrado")
		}
		return nil, err
	}

	existente.EstadoPago = pago.EstadoPago
	existente.MontoPagado = pago.MontoPagado

	if err := database.DB.Save(&existente).Error; err != nil {
		return nil, err
	}

	return &existente, nil
}

func (s *PagoService) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Pago{}, id).Error; err != nil {
		return err
	}
	return nil
}
