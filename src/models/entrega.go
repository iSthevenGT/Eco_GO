package models

type EstadoEntrega string

const (
	EntregaPreparacion EstadoEntrega = "preparacion"
	EntregaEnviada     EstadoEntrega = "enviada"
	EntregaEntregada   EstadoEntrega = "entregada"
	EntregaCancelada   EstadoEntrega = "cancelada"
)

type Entrega struct {
	ID            uint          `json:"id" gorm:"primaryKey"`
	NumeroEntrega int           `json:"numero_entrega" gorm:"not null"`
	EstadoEntrega EstadoEntrega `json:"estado_entrega" gorm:"default:preparacion"`

	// Relaciones
	OrdenID      uint               `json:"orden_id" gorm:"not null;uniqueIndex"`
	Orden        *Orden             `json:"orden,omitempty" gorm:"foreignKey:OrdenID"`
	RepartidorID uint               `json:"repartidor_id" gorm:"not null"`
	Repartidor   *Repartidor        `json:"repartidor,omitempty" gorm:"foreignKey:RepartidorID"`
	Direcciones  []EntregaDireccion `json:"direcciones,omitempty" gorm:"foreignKey:EntregaID"`
}

type TipoDireccion string

const (
	DireccionParada  TipoDireccion = "parada"
	DireccionEntrega TipoDireccion = "entrega"
)

type EntregaDireccion struct {
	ID   uint          `json:"id" gorm:"primaryKey"`
	Tipo TipoDireccion `json:"tipo" gorm:"not null"`

	// Relaciones
	EntregaID   uint       `json:"entrega_id" gorm:"not null"`
	Entrega     *Entrega   `json:"entrega,omitempty" gorm:"foreignKey:EntregaID"`
	DireccionID uint       `json:"direccion_id" gorm:"not null"`
	Direccion   *Direccion `json:"direccion,omitempty" gorm:"foreignKey:DireccionID"`
}
