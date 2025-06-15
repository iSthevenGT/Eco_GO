package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/models"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ComercianteService struct{}

func NewComercianteService() *ComercianteService {
	return &ComercianteService{}
}

func (s *ComercianteService) ObtenerTodos() ([]models.Comerciante, error) {
	var comerciantes []models.Comerciante
	if err := database.DB.Preload("Telefono").Find(&comerciantes).Error; err != nil {
		return nil, err
	}
	return comerciantes, nil
}

func (s *ComercianteService) ObtenerPorID(id uint) (*models.Comerciante, error) {
	var comerciante models.Comerciante
	if err := database.DB.Preload("Telefono").Preload("Productos").First(&comerciante, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comerciante no encontrado")
		}
		return nil, err
	}
	return &comerciante, nil
}

func (s *ComercianteService) Eliminar(id uint) error {
	if err := database.DB.Delete(&models.Comerciante{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *ComercianteService) Actualizar(id uint, comerciante models.Comerciante) (*models.Comerciante, error) {
	var existente models.Comerciante
	if err := database.DB.First(&existente, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comerciante no encontrado")
		}
		return nil, err
	}

	existente.Nombre = comerciante.Nombre
	if err := database.DB.Save(&existente).Error; err != nil {
		return nil, err
	}

	return &existente, nil
}

func (s *ComercianteService) ObtenerProductos(id uint) ([]models.Producto, error) {
	var comerciante models.Comerciante
	if err := database.DB.Preload("Productos").First(&comerciante, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comerciante no encontrado")
		}
		return nil, err
	}
	return comerciante.Productos, nil
}

func (s *ComercianteService) ObtenerProductoPorID(comercianteID, productoID uint) (*models.Producto, error) {
	var producto models.Producto
	if err := database.DB.Where("id = ? AND comerciante_id = ?", productoID, comercianteID).First(&producto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("producto no encontrado")
		}
		return nil, err
	}
	return &producto, nil
}

func (s *ComercianteService) CompletarRegistro(id uint, nit string, camaraComercio, rut *multipart.FileHeader) (*models.Comerciante, error) {
	var comerciante models.Comerciante
	if err := database.DB.First(&comerciante, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("comerciante no encontrado")
		}
		return nil, err
	}

	comerciante.NIT = nit

	// Guardar cámara de comercio
	if camaraComercio != nil {
		ccPath, err := s.guardarArchivo(camaraComercio, "camaraComercio", id)
		if err != nil {
			return nil, fmt.Errorf("error al guardar cámara de comercio: %v", err)
		}
		comerciante.CamaraComercio = ccPath
	}

	// Guardar RUT
	if rut != nil {
		rutPath, err := s.guardarArchivo(rut, "rut", id)
		if err != nil {
			return nil, fmt.Errorf("error al guardar RUT: %v", err)
		}
		comerciante.RUT = rutPath
	}

	if err := database.DB.Save(&comerciante).Error; err != nil {
		return nil, err
	}

	return &comerciante, nil
}

func (s *ComercianteService) guardarArchivo(fileHeader *multipart.FileHeader, tipo string, id uint) (string, error) {
	// Crear directorio si no existe
	uploadDir := "uploads/usuarios/documentos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	// Generar nombre único
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s_%s%s", uuid.New().String(), tipo, ext)
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

	return fmt.Sprintf("%s/usuarios/documentos/%s", baseURL, filename), nil
}
