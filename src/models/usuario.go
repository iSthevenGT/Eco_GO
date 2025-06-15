package models

import (
	"time"

	"gorm.io/gorm"
)

type Usuario struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Imagen     string         `json:"imagen" gorm:"type:text"`
	Nombre     string         `json:"nombre" gorm:"type:varchar(255);not null"`
	Email      string         `json:"Email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Contrasena string         `json:"-" gorm:"type:varchar(255);not null"`
	Rol        string         `json:"rol" gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`

	// Relación con teléfono (1:1)
	Telefono *Telefono `json:"telefono,omitempty" gorm:"foreignKey:UsuarioID"`

	// OPCIÓN 1: Relaciones opcionales - solo una será válida según el rol
	Consumidor  *Consumidor  `json:"consumidor,omitempty" gorm:"foreignKey:UsuarioID"`
	Comerciante *Comerciante `json:"comerciante,omitempty" gorm:"foreignKey:UsuarioID"`
	Repartidor  *Repartidor  `json:"repartidor,omitempty" gorm:"foreignKey:UsuarioID"`
}

// OPCIÓN 1: Con clave foránea separada (MÁS SEGURO)
type Consumidor struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	UsuarioID uint `json:"usuario_id" gorm:"uniqueIndex;not null"`

	// Relación inversa
	Usuario Usuario `json:"usuario" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Puntos      int                `json:"puntos" gorm:"default:0"`
	Direcciones []UsuarioDireccion `json:"direcciones,omitempty" gorm:"foreignKey:UsuarioID"`
	Ordenes     []Orden            `json:"ordenes,omitempty" gorm:"foreignKey:ConsumidorID"`
}

type Comerciante struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	UsuarioID uint `json:"usuario_id" gorm:"uniqueIndex;not null"`

	// Relación inversa
	Usuario Usuario `json:"usuario" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	NIT            string     `json:"nit" gorm:"type:varchar(100)"`
	RUT            string     `json:"rut" gorm:"type:varchar(100)"`
	CamaraComercio string     `json:"camara_comercio" gorm:"type:varchar(100)"`
	Productos      []Producto `json:"productos,omitempty" gorm:"foreignKey:ComercianteID"`
	Direccion      *Direccion `json:"direccion,omitempty" gorm:"foreignKey:DireccionID"`
	DireccionID    *uint      `json:"direccion_id"`
}

type Repartidor struct {
	ID        uint `json:"id" gorm:"primaryKey"`
	UsuarioID uint `json:"usuario_id" gorm:"uniqueIndex;not null"`

	// Relación inversa
	Usuario Usuario `json:"usuario" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Placa    string `json:"placa" gorm:"type:varchar(20)"`
	SOAT     string `json:"soat" gorm:"type:varchar(100)"`
	Licencia string `json:"licencia" gorm:"type:varchar(100)"`
	Tecno    string `json:"tecno" gorm:"type:varchar(100)"`
}
