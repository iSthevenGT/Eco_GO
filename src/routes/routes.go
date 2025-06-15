package routes

import (
	"Eco_GO/src/handlers"
	"Eco_GO/src/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Servir archivos estáticos
	router.Static("/productos", "./uploads/productos")
	router.Static("/usuarios", "./uploads/usuarios")
	router.Static("/usuarios/documentos", "./uploads/usuarios/documentos")

	// Inicializar controladores
	authHandler := handlers.NewAuthHandler()
	adminHandler := handlers.NewAdminHandler()
	consumidorHandler := handlers.NewConsumidorHandler()
	comercianteHandler := handlers.NewComercianteHandler()

	// Rutas públicas
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/validate-token", authHandler.ValidateToken)
	}

	// Aplicar middleware de autenticación para rutas protegidas
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Rutas de administrador
		admin := api.Group("/admin")
		admin.Use(middleware.RequireRole("ADMIN")) // Asumiendo que existe un rol ADMIN
		{
			admin.GET("/consumidores/todos", adminHandler.ObtenerConsumidores)
			admin.GET("/consumidores/:idConsumidor", adminHandler.ObtenerConsumidorPorID)
			admin.GET("/comerciantes/todos", adminHandler.ObtenerComerciantes)
			admin.GET("/comerciantes/:idComerciante", adminHandler.ObtenerComerciantePorID)
		}

		// Rutas de consumidores
		consumidores := api.Group("/consumidores")
		consumidores.Use(middleware.RequireRole("CONSUMIDOR"))
		{
			consumidores.POST("/:id/establecerImagen", consumidorHandler.EstablecerImagen)
			consumidores.POST("/:id/crearTelefono", consumidorHandler.CrearTelefono)
			consumidores.POST("/:id/crearDireccion", consumidorHandler.CrearDireccion)

			// Productos
			consumidores.GET("/:id/productos/todos", consumidorHandler.ObtenerProductos)
			consumidores.GET("/:id/productos/:idProducto", consumidorHandler.ObtenerProductoPorID)

			// Carrito
			consumidores.POST("/:id/productos/:idProducto/agregar", consumidorHandler.AgregarAlCarrito)
			consumidores.GET("/:id/carrito", consumidorHandler.VerCarrito)
			consumidores.GET("/:id/carrito/:productoId/eliminar", consumidorHandler.EliminarProductoCarrito)
			consumidores.GET("/:id/carrito/:idProducto/cambiarCantidad", consumidorHandler.CambiarCantidadProducto)
			consumidores.GET("/:id/carrito/limpiar", consumidorHandler.LimpiarCarrito)

			// Órdenes
			consumidores.POST("/:id/carrito/ordenar", consumidorHandler.CrearOrden)
			consumidores.GET("/:id/ordenes/:idOrden", consumidorHandler.ObtenerOrden)
			consumidores.POST("/:id/ordenes/:idOrden/cancelar", consumidorHandler.CancelarOrden)
		}

		// Rutas de comerciantes
		comerciantes := api.Group("/comerciantes")
		comerciantes.Use(middleware.RequireRole("COMERCIANTE"))
		{
			comerciantes.POST("/:idComerciante/establecerImagen", comercianteHandler.EstablecerImagen)
			comerciantes.POST("/:idComerciante/crearTelefono", comercianteHandler.CrearTelefono)

			// Productos
			comerciantes.POST("/:idComerciante/crearProducto", comercianteHandler.CrearProducto)
			comerciantes.GET("/:idComerciante/productos/todos", comercianteHandler.ObtenerProductos)
			comerciantes.GET("/:idComerciante/productos/:idProducto", comercianteHandler.ObtenerProductoPorID)
			comerciantes.POST("/:idComerciante/actualizarProducto/:idProducto", comercianteHandler.ActualizarProducto)

			// Órdenes
			comerciantes.GET("/:idComerciante/ordenes/todos", comercianteHandler.ObtenerOrdenes)
			comerciantes.GET("/:idComerciante/ordenes/:idOrden", comercianteHandler.ObtenerOrden)
			comerciantes.POST("/:idComerciante/ordenes/:idOrden/confirmar", comercianteHandler.ConfirmarOrden)

			// Preparación de órdenes
			comerciantes.GET("/:idComerciante/ordenes/preparacion", comercianteHandler.OrdenesPreparacion)
			comerciantes.GET("/:idComerciante/ordenes/preparacion/:idOrden", comercianteHandler.OrdenPreparacion)

			// Registro completo
			comerciantes.POST("/:id/completarRegistro", comercianteHandler.CompletarRegistro)
		}

		// Rutas de repartidores
		repartidores := api.Group("/repartidores")
		repartidores.Use(middleware.RequireRole("REPARTIDOR"))
		{
			// Aquí irían las rutas específicas de repartidores cuando las implementes
		}
	}
}
