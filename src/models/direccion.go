package models

import (
	"time"

	"gorm.io/gorm"
)

type Direccion struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Nombre    string         `json:"nombre" gorm:"not null"`
	Domicilio string         `json:"domicilio" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UsuarioDireccion struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	UsuarioID   uint `json:"usuario_id" gorm:"not null"`
	DireccionID uint `json:"direccion_id" gorm:"not null"`

	// Relaciones
	Usuario   *Usuario   `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
	Direccion *Direccion `json:"direccion,omitempty" gorm:"foreignKey:DireccionID"`
}

type Telefono struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Indicativo string `json:"indicativo" gorm:"not null"`
	Numero     string `json:"numero" gorm:"not null"`
	UsuarioID  uint   `json:"usuario_id" gorm:"not null;uniqueIndex"`

	// Relaciones
	Usuario *Usuario `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}
