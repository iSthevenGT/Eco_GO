package handlers

import (
	"Eco_GO/src/models"
	"Eco_GO/src/services"
	"Eco_GO/src/src/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConsumidorHandler struct {
	carritoService *services.CarritoService
}

func NewConsumidorHandler() *ConsumidorHandler {
	return &ConsumidorHandler{
		carritoService: services.NewCarritoService(),
	}
}

func (ctrl *ConsumidorHandler) ObtenerProductos(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Verificar que el consumidor existe
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, idConsumidor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consumidor no encontrado"})
		return
	}

	var productos []models.Producto
	if err := database.DB.Preload("Comerciante").Find(&productos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}

func (ctrl *ConsumidorHandler) ObtenerProductoPorID(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	idProducto, _ := strconv.ParseUint(c.Param("productId"), 10, 32)

	// Verificar que el consumidor existe
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, idConsumidor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consumidor no encontrado"})
		return
	}

	var producto models.Producto
	if err := database.DB.Preload("Comerciante").First(&producto, idProducto).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, producto)
}

func (ctrl *ConsumidorHandler) AgregarAlCarrito(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	idProducto, _ := strconv.ParseUint(c.Param("productId"), 10, 32)

	var req struct {
		Cantidad int `json:"cantidad" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.carritoService.AgregarProducto(uint(idConsumidor), uint(idProducto), req.Cantidad)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto agregado al carrito"})
}

func (ctrl *ConsumidorHandler) VerCarrito(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	carrito := ctrl.carritoService.ObtenerCarrito(uint(idConsumidor))
	c.JSON(http.StatusOK, carrito)
}

func (ctrl *ConsumidorHandler) EliminarProductoCarrito(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	idProducto, _ := strconv.ParseUint(c.Param("productId"), 10, 32)

	err := ctrl.carritoService.EliminarProducto(uint(idConsumidor), uint(idProducto))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	carrito := ctrl.carritoService.ObtenerCarrito(uint(idConsumidor))
	c.JSON(http.StatusOK, carrito)
}

func (ctrl *ConsumidorHandler) LimpiarCarrito(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	ctrl.carritoService.LimpiarCarrito(uint(idConsumidor))
	c.JSON(http.StatusOK, gin.H{"message": "Carrito limpiado con éxito"})
}

func (ctrl *ConsumidorHandler) CrearDireccion(c *gin.Context) {
	idConsumidor, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var direccion models.Direccion
	if err := c.ShouldBindJSON(&direccion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el consumidor existe
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, idConsumidor).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Consumidor no encontrado"})
		return
	}

	// Crear la dirección
	if err := database.DB.Create(&direccion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crear la relación UsuarioDireccion
	usuarioDireccion := models.UsuarioDireccion{
		UsuarioID:   uint(idConsumidor),
		DireccionID: direccion.ID,
	}

	if err := database.DB.Create(&usuarioDireccion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, usuarioDireccion)
}
