package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductoService struct{}

func NewProductoService() *ProductoService {
	return &ProductoService{}
}

func (s *ProductoService) ObtenerTodos(consumidorID uint) ([]models.Producto, error) {
	// Verificar que el consumidor existe
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, consumidorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("consumidor no encontrado")
		}
		return nil, err
	}

	var productos []models.Producto
	if err := database.DB.Preload("Comerciante").Find(&productos).Error; err != nil {
		return nil, err
	}
	return productos, nil
}

func (s *ProductoService) ObtenerPorID(consumidorID, productoID uint) (*models.Producto, error) {
	// Verificar que el consumidor existe
	var consumidor models.Consumidor
	if err := database.DB.First(&consumidor, consumidorID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("consumidor no encontrado")
		}
		return nil, err
	}

	var producto models.Producto
	if err := database.DB.Preload("Comerciante").First(&producto, productoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}
	return &producto, nil
}

func (s *ProductoService) Crear(comercianteID uint, productoJSON string, imagen *multipart.FileHeader) (*models.Producto, error) {
	// Verificar que el comerciante existe
	var comerciante models.Comerciante
	if err := database.DB.First(&comerciante, comercianteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comerciante no encontrado")
		}
		return nil, err
	}

	// Parsear JSON del producto
	var productoData map[string]interface{}
	if err := json.Unmarshal([]byte(productoJSON), &productoData); err != nil {
		return nil, fmt.Errorf("error al parsear JSON del producto: %v", err)
	}

	// Crear producto
	producto := models.Producto{
		Nombre:        productoData["nombre"].(string),
		Descripcion:   getString(productoData, "descripcion"),
		Tipo:          getString(productoData, "tipo"),
		Precio:        int(getFloat64(productoData, "precio")),
		Stock:         int(getFloat64(productoData, "stock")),
		ComercianteID: comercianteID,
	}

	// Guardar imagen si se proporciona
	if imagen != nil {
		imagenURL, err := s.guardarImagen(imagen)
		if err != nil {
			return nil, fmt.Errorf("error al guardar imagen: %v", err)
		}
		producto.Imagen = imagenURL
	}

	if err := database.DB.Create(&producto).Error; err != nil {
		return nil, err
	}

	return &producto, nil
}

func (s *ProductoService) Actualizar(comercianteID, productoID uint, productoJSON string, imagen *multipart.FileHeader) (*models.Producto, error) {
	// Verificar que el producto existe y pertenece al comerciante
	var producto models.Producto
	if err := database.DB.Where("id = ? AND comerciante_id = ?", productoID, comercianteID).First(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}

	// Parsear JSON del producto
	var productoData map[string]interface{}
	if err := json.Unmarshal([]byte(productoJSON), &productoData); err != nil {
		return nil, fmt.Errorf("error al parsear JSON del producto: %v", err)
	}

	// Actualizar campos
	producto.Nombre = productoData["nombre"].(string)
	producto.Descripcion = getString(productoData, "descripcion")
	producto.Tipo = getString(productoData, "tipo")
	producto.Precio = int(getFloat64(productoData, "precio"))
	producto.Stock = int(getFloat64(productoData, "stock"))

	// Actualizar imagen si se proporciona
	if imagen != nil {
		imagenURL, err := s.guardarImagen(imagen)
		if err != nil {
			return nil, fmt.Errorf("error al guardar imagen: %v", err)
		}
		producto.Imagen = imagenURL
	}

	if err := database.DB.Save(&producto).Error; err != nil {
		return nil, err
	}

	return &producto, nil
}

func (s *ProductoService) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Producto{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *ProductoService) guardarImagen(fileHeader *multipart.FileHeader) (string, error) {
	// Crear directorio si no existe
	uploadDir := "uploads/productos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Generar nombre único
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Abrir archivo desde el header
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Crear archivo de destino
	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copiar contenido
	if _, err := dst.ReadFrom(file); err != nil {
		return "", err
	}

	// Retornar URL completa
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return fmt.Sprintf("%s/productos/%s", baseURL, filename), nil
}

// Funciones auxiliares para extraer valores del map
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getFloat64(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		}
	}
	return 0
}
