package models

import (
	"time"

	"gorm.io/gorm"
)

type Producto struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Imagen      string         `json:"imagen"`
	Nombre      string         `json:"nombre" gorm:"not null"`
	Precio      int            `json:"precio" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	Tipo        string         `json:"tipo"`
	Descripcion string         `json:"descripcion" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	ComercianteID uint         `json:"comerciante_id" gorm:"not null"`
	Comerciante   *Comerciante `json:"comerciante,omitempty" gorm:"foreignKey:ComercianteID"`
	Puntuaciones  []Puntuacion `json:"puntuaciones,omitempty" gorm:"foreignKey:ProductoID"`
}

type Puntuacion struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Puntos    int       `json:"puntos" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relaciones
	UsuarioID  uint      `json:"usuario_id" gorm:"not null"`
	Usuario    *Usuario  `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
	ProductoID uint      `json:"producto_id" gorm:"not null"`
	Producto   *Producto `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
}
