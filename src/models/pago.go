package models

import "time"

type EstadoPago string
type MetodoPago string

const (
	PagoPendiente   EstadoPago = "pendiente"
	PagoAprobado    EstadoPago = "aprobado"
	PagoRechazado   EstadoPago = "rechazado"
	PagoReembolsado EstadoPago = "reembolsado"

	PagoCredito  MetodoPago = "credito"
	PagoPSE      MetodoPago = "pse"
	PagoEfectivo MetodoPago = "efectivo"
)

type Pago struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	EstadoPago  EstadoPago `json:"estado_pago" gorm:"default:pendiente"`
	MetodoPago  MetodoPago `json:"metodo_pago" gorm:"not null"`
	MontoPagado int        `json:"monto_pagado" gorm:"not null"`
	CreatedAt   time.Time  `json:"created_at"`

	// Relaciones
	OrdenID     uint   `json:"orden_id" gorm:"not null;uniqueIndex"`
	Orden       *Orden `json:"orden,omitempty" gorm:"foreignKey:OrdenID"`
	FechaPagoID uint   `json:"fecha_pago_id" gorm:"not null"`
	FechaPago   *Fecha `json:"fecha_pago,omitempty" gorm:"foreignKey:FechaPagoID"`
}
