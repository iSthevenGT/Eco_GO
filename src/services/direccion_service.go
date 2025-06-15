package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
)

type DireccionService struct{}

func NewDireccionService() *DireccionService {
	return &DireccionService{}
}

func (s *DireccionService) Crear(direccion models.Direccion) (*models.Direccion, error) {
	if err := database.DB.Create(&direccion).Error; err != nil {
		return nil, err
	}
	return &direccion, nil
}

func (s *DireccionService) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Direccion{}, id).Error; err != nil {
		return err
	}
	return nil
}
