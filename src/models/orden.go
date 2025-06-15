package models

import (
	"time"

	"gorm.io/gorm"
)

type EstadoOrden string

const (
	EstadoPendiente   EstadoOrden = "PENDIENTE"
	EstadoConfirmada  EstadoOrden = "CONFIRMADA"
	EstadoCancelada   EstadoOrden = "CANCELADA"
	EstadoReembolsada EstadoOrden = "REEMBOLSADA"
)

type Orden struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MontoTotal  int            `json:"monto_total" gorm:"not null"`
	EstadoOrden EstadoOrden    `json:"estado_orden" gorm:"default:PENDIENTE"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relaciones
	ConsumidorID       uint            `json:"consumidor_id" gorm:"not null"`
	Consumidor         *Consumidor     `json:"consumidor,omitempty" gorm:"foreignKey:ConsumidorID"`
	FechaOrdenID       uint            `json:"fecha_orden_id" gorm:"not null"`
	FechaOrden         *Fecha          `json:"fecha_orden,omitempty" gorm:"foreignKey:FechaOrdenID"`
	DireccionEntregaID *uint           `json:"direccion_entrega_id"`
	DireccionEntrega   *Direccion      `json:"direccion_entrega,omitempty" gorm:"foreignKey:DireccionEntregaID"`
	Productos          []OrdenProducto `json:"productos,omitempty" gorm:"foreignKey:OrdenID"`
	Pago               *Pago           `json:"pago,omitempty" gorm:"foreignKey:OrdenID"`
}

type OrdenProducto struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	Cantidad int  `json:"cantidad" gorm:"default:1"`

	// Relaciones
	OrdenID    uint      `json:"orden_id" gorm:"not null"`
	Orden      *Orden    `json:"orden,omitempty" gorm:"foreignKey:OrdenID"`
	ProductoID uint      `json:"producto_id" gorm:"not null"`
	Producto   *Producto `json:"producto,omitempty" gorm:"foreignKey:ProductoID"`
}

type Fecha struct {
	ID   uint `json:"id" gorm:"primaryKey"`
	Anio int  `json:"anio" gorm:"not null"`
	Mes  int  `json:"mes" gorm:"not null"`
	Dia  int  `json:"dia" gorm:"not null"`
}
