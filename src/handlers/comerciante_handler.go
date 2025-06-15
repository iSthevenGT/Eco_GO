package handlers

import (
	"Eco_GO/src/models"
	"Eco_GO/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ComercianteHandler struct {
	comercianteService        *services.ComercianteService
	productoService           *services.ProductoService
	ordenService              *services.OrdenService
	telefonoService           *services.TelefonoService
	usuarioService            *services.UsuarioService
	preparacionOrdenesService *services.PreparacionOrdenesService
}

func NewComercianteHandler() *ComercianteHandler {
	return &ComercianteHandler{
		comercianteService:        services.NewComercianteService(),
		productoService:           services.NewProductoService(),
		ordenService:              services.NewOrdenService(),
		telefonoService:           services.NewTelefonoService(),
		usuarioService:            services.NewUsuarioService(),
		preparacionOrdenesService: services.NewPreparacionOrdenesService(),
	}
}

func (ctrl *ComercianteHandler) EstablecerImagen(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	imagen, err := c.FormFile("imagen")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagen requerida"})
		return
	}

	authHeader := c.GetHeader("Authorization")
	token := authHeader[7:] // Remover "Bearer "

	usuario, err := ctrl.usuarioService.SetImagen(uint(idComerciante), imagen, token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

func (ctrl *ComercianteHandler) CrearTelefono(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	var telefono models.Telefono
	if err := c.ShouldBindJSON(&telefono); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	telefonoCreado, err := ctrl.telefonoService.Crear(uint(idComerciante), telefono)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, telefonoCreado)
}

func (ctrl *ComercianteHandler) CrearProducto(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	productoJSON := c.PostForm("producto")
	imagen, err := c.FormFile("imagen")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagen requerida"})
		return
	}

	producto, err := ctrl.productoService.Crear(uint(idComerciante), productoJSON, imagen)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, producto)
}

func (ctrl *ComercianteHandler) ObtenerProductos(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	productos, err := ctrl.comercianteService.ObtenerProductos(uint(idComerciante))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}

func (ctrl *ComercianteHandler) ObtenerProductoPorID(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)
	idProducto, _ := strconv.ParseUint(c.Param("idProducto"), 10, 32)

	producto, err := ctrl.comercianteService.ObtenerProductoPorID(uint(idComerciante), uint(idProducto))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, producto)
}

func (ctrl *ComercianteHandler) ActualizarProducto(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)
	idProducto, _ := strconv.ParseUint(c.Param("idProducto"), 10, 32)

	productoJSON := c.PostForm("producto")
	imagen, err := c.FormFile("imagen")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Imagen requerida"})
		return
	}

	producto, err := ctrl.productoService.Actualizar(uint(idComerciante), uint(idProducto), productoJSON, imagen)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, producto)
}

func (ctrl *ComercianteHandler) ObtenerOrdenes(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	ordenes, err := ctrl.ordenService.ObtenerTodosPorComerciante(uint(idComerciante))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ordenes)
}

func (ctrl *ComercianteHandler) ObtenerOrden(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)
	idOrden, _ := strconv.ParseUint(c.Param("idOrden"), 10, 32)

	productos, err := ctrl.ordenService.ProductosPorIDAndComerciante(uint(idComerciante), uint(idOrden))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}

func (ctrl *ComercianteHandler) ConfirmarOrden(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)
	idOrden, _ := strconv.ParseUint(c.Param("idOrden"), 10, 32)

	orden, err := ctrl.preparacionOrdenesService.AgregarOrden(uint(idComerciante), uint(idOrden))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orden)
}

func (ctrl *ComercianteHandler) OrdenesPreparacion(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	ordenes := ctrl.preparacionOrdenesService.ObtenerOrdenes(uint(idComerciante))
	c.JSON(http.StatusOK, ordenes)
}

func (ctrl *ComercianteHandler) OrdenPreparacion(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)
	idOrden, _ := strconv.ParseUint(c.Param("idOrden"), 10, 32)

	orden, err := ctrl.preparacionOrdenesService.ObtenerOrden(uint(idComerciante), uint(idOrden))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orden)
}

func (ctrl *ComercianteHandler) CompletarRegistro(c *gin.Context) {
	idComerciante, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	nit := c.PostForm("nit")
	camaraComercio, err1 := c.FormFile("camaraComercio")
	rut, err2 := c.FormFile("rut")

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivos requeridos"})
		return
	}

	comerciante, err := ctrl.comercianteService.CompletarRegistro(uint(idComerciante), nit, camaraComercio, rut)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comerciante)
}
