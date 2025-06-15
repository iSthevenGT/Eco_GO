package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"

	"gorm.io/gorm"
)

type UsuarioDireccionService struct {
	direccionService *DireccionService
}

func NewUsuarioDireccionService() *UsuarioDireccionService {
	return &UsuarioDireccionService{
		direccionService: NewDireccionService(),
	}
}

func (s *UsuarioDireccionService) Crear(usuarioID uint, direccion models.Direccion) (*models.UsuarioDireccion, error) {
	// Verificar que el usuario existe
	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Crear la dirección
	nuevaDireccion, err := s.direccionService.Crear(direccion)
	if err != nil {
		return nil, err
	}

	// Crear la relación usuario-dirección
	usuarioDireccion := models.UsuarioDireccion{
		UsuarioID:   usuarioID,
		DireccionID: nuevaDireccion.ID,
	}

	if err := database.DB.Create(&usuarioDireccion).Error; err != nil {
		return nil, err
	}

	// Cargar relaciones
	if err := database.DB.Preload("Usuario").Preload("Direccion").First(&usuarioDireccion, usuarioDireccion.ID).Error; err != nil {
		return nil, err
	}

	return &usuarioDireccion, nil
}
