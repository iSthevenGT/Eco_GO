package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"

	"gorm.io/gorm"
)

type TelefonoService struct{}

func NewTelefonoService() *TelefonoService {
	return &TelefonoService{}
}

func (s *TelefonoService) ObtenerTodos() ([]models.Telefono, error) {
	var telefonos []models.Telefono
	if err := database.DB.Find(&telefonos).Error; err != nil {
		return nil, err
	}
	return telefonos, nil
}

func (s *TelefonoService) ObtenerPorID(id uint) (*models.Telefono, error) {
	var telefono models.Telefono
	if err := database.DB.First(&telefono, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teléfono no encontrado")
		}
		return nil, err
	}
	return &telefono, nil
}

func (s *TelefonoService) Crear(usuarioID uint, telefono models.Telefono) (*models.Telefono, error) {
	// Verificar que el usuario existe
	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Verificar que el usuario no tenga ya un teléfono
	var existingPhone models.Telefono
	if err := database.DB.Where("usuario_id = ?", usuarioID).First(&existingPhone).Error; err == nil {
		return nil, errors.New("el usuario ya tiene un teléfono asignado")
	}

	telefono.UsuarioID = usuarioID
	if err := database.DB.Create(&telefono).Error; err != nil {
		return nil, err
	}

	return &telefono, nil
}

func (s *TelefonoService) Actualizar(telefono models.Telefono) (*models.Telefono, error) {
	if err := database.DB.Save(&telefono).Error; err != nil {
		return nil, err
	}
	return &telefono, nil
}

func (s *TelefonoService) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Telefono{}, id).Error; err != nil {
		return err
	}
	return nil
}
