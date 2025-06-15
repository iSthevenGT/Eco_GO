package repositories

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
}

func NewUsuarioRepository() *UsuarioRepository {
	return &UsuarioRepository{}
}

func (rep *UsuarioRepository) Crear(usuario models.Usuario) (*models.Usuario, error) {
	if err := database.DB.Create(&usuario).Error; err != nil {
		return nil, err
	}
	return &usuario, nil

}

func (rep *UsuarioRepository) ObtenerPorID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	if err := database.DB.First(&usuario, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}

func (rep *UsuarioRepository) Actualizar(usuario models.Usuario) (*models.Usuario, error) {
	var existente models.Usuario
	if err := database.DB.First(&existente, usuario.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	existente.Nombre = usuario.Nombre
	existente.Email = usuario.Email

	if err := database.DB.Save(&existente).Error; err != nil {
		return nil, err
	}

	return &existente, nil
}

func (rep *UsuarioRepository) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Usuario{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (rep *UsuarioRepository) ObtenerTodos() ([]models.Usuario, error) {
	var usuarios []models.Usuario
	if err := database.DB.Find(&usuarios).Error; err != nil {
		return nil, err
	}
	return usuarios, nil
}

func (rep *UsuarioRepository) ObtenerPorEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	if err := database.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &usuario, nil
}
