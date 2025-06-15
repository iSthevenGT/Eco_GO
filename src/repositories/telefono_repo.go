package repositories

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"

	"gorm.io/gorm"
)

type TelefonoRepository struct{}

func NewTelefonoRepository() *TelefonoRepository {
	return &TelefonoRepository{}
}

// Método principal que puede usarse independientemente
func (rep *TelefonoRepository) Crear(usuarioID uint, telefono map[string]string) (*models.Telefono, error) {
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	telefonoCreado, err := rep.CrearConTx(tx, usuarioID, telefono)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return telefonoCreado, nil
}

// Método auxiliar que usa una transacción existente
func (rep *TelefonoRepository) CrearConTx(tx *gorm.DB, usuarioID uint, telefono map[string]string) (*models.Telefono, error) {
	// Validar que se proporcionen datos de teléfono
	if telefono == nil || telefono["numero"] == "" || telefono["indicativo"] == "" {
		return nil, nil // No crear teléfono si no hay datos válidos
	}

	telefonoCreado := models.Telefono{
		Numero:     telefono["numero"],
		Indicativo: telefono["indicativo"],
		UsuarioID:  usuarioID,
	}

	if err := tx.Create(&telefonoCreado).Error; err != nil {
		return nil, err
	}

	return &telefonoCreado, nil
}
