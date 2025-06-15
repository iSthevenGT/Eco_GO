package repositories

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	telefonoRepo *TelefonoRepository
}

func NewUsuarioRepository() *UsuarioRepository {
	return &UsuarioRepository{
		telefonoRepo: NewTelefonoRepository(),
	}
}

func (rep *UsuarioRepository) Crear(usuario *models.Usuario, telefono map[string]string) (*models.Usuario, error) {
	// Iniciar transacción para garantizar consistencia
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(usuario).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Crear entidad relacionada según el rol
	switch usuario.Rol {
	case "CONSUMIDOR":
		consumidor := models.Consumidor{
			UsuarioID: usuario.ID,
			Puntos:    0,
		}
		if err := tx.Create(&consumidor).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	case "COMERCIANTE":
		comerciante := models.Comerciante{
			UsuarioID: usuario.ID,
		}
		if err := tx.Create(&comerciante).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	case "REPARTIDOR":
		repartidor := models.Repartidor{
			UsuarioID: usuario.ID,
		}
		if err := tx.Create(&repartidor).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

	default:
		tx.Rollback()
		return nil, errors.New("rol no válido")
	}

	// Crear el teléfono usando la MISMA transacción
	telefonoCreado, err := rep.telefonoRepo.CrearConTx(tx, usuario.ID, telefono)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	usuario.Telefono = telefonoCreado

	// Confirmar transacción
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return usuario, nil
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

func (rep *UsuarioRepository) ValidarEmail(email string) error {
	var usuario models.Usuario
	if err := database.DB.Where("email = ?", email).First(&usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // Email disponible
		}
		return err // Error de base de datos
	}
	return errors.New("el Email ya está en uso") // Email ya existe
}
