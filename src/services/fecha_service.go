package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
)

type FechaService struct{}

func NewFechaService() *FechaService {
	return &FechaService{}
}

func (s *FechaService) Crear(fecha models.Fecha) (*models.Fecha, error) {
	if err := database.DB.Create(&fecha).Error; err != nil {
		return nil, err
	}
	return &fecha, nil
}

func (s *FechaService) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Fecha{}, id).Error; err != nil {
		return err
	}
	return nil
}
