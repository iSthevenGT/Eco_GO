package handlers

import (
	"Eco_GO/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminHanler struct {
	consumidorService  *services.ConsumidorService
	comercianteService *services.ComercianteService
}

func NewAdminHandler() *AdminHanler {
	return &AdminHanler{
		consumidorService:  services.NewConsumidorService(),
		comercianteService: services.NewComercianteService(),
	}
}

func (ctrl *AdminHanler) ObtenerConsumidores(c *gin.Context) {
	consumidores, err := ctrl.consumidorService.ObtenerTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consumidores)
}

func (ctrl *AdminHanler) ObtenerConsumidorPorID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("idConsumidor"), 10, 32)

	consumidor, err := ctrl.consumidorService.ObtenerPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, consumidor)
}

func (ctrl *AdminHanler) ObtenerComerciantes(c *gin.Context) {
	comerciantes, err := ctrl.comercianteService.ObtenerTodos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comerciantes)
}

func (ctrl *AdminHanler) ObtenerComerciantePorID(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("idComerciante"), 10, 32)

	comerciante, err := ctrl.comercianteService.ObtenerPorID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comerciante)
}
