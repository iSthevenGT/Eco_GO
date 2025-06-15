package dto

type ProductoDTO struct {
	ID       uint   `json:"id"`
	Imagen   string `json:"imagen"`
	Nombre   string `json:"nombre"`
	Precio   int    `json:"precio"`
	Cantidad int    `json:"cantidad"`
}

type CarritoDTO struct {
	Productos []ProductoDTO `json:"productos"`
	Total     int           `json:"total"`
}

type OrdenesPrepDTO struct {
	Ordenes []interface{} `json:"ordenes"` // Usar interface{} para flexibilidad
}

type LoginRequest struct {
	Email      string `json:"Email" binding:"required,email"`
	Contrasena string `json:"contrasena" binding:"required"`
}

type RegisterRequest struct {
	Nombre     string            `json:"nombre" binding:"required"`
	Email      string            `json:"Email" binding:"required,email"`
	Contrasena string            `json:"contrasena" binding:"required,min=4"`
	Rol        string            `json:"rol" binding:"required"`
	Telefono   map[string]string `json:"telefono"`
}

type CreateProductoRequest struct {
	Nombre      string `json:"nombre" binding:"required"`
	Descripcion string `json:"descripcion"`
	Tipo        string `json:"tipo"`
	Precio      int    `json:"precio" binding:"required,min=0"`
	Stock       int    `json:"stock" binding:"required,min=0"`
}

type CreateOrdenRequest struct {
	FechaOrden       map[string]int         `json:"fechaOrden" binding:"required"`
	DireccionEntrega map[string]interface{} `json:"direccionEntrega"`
	Pago             map[string]string      `json:"pago" binding:"required"`
}
