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

type UsuarioService struct{}

func NewUsuarioService() *UsuarioService {
	return &UsuarioService{}
}

func (s *UsuarioService) SetImagen(usuarioID uint, imagen *multipart.FileHeader, token string) (*models.Usuario, error) {
	// Aquí deberías validar el token y verificar permisos
	// Por simplicidad, omitimos esa validación en este ejemplo

	var usuario models.Usuario
	if err := database.DB.First(&usuario, usuarioID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	// Eliminar imagen anterior si existe
	if usuario.Imagen != "" {
		// Aquí podrías implementar la lógica para eliminar el archivo anterior
	}

	// Guardar nueva imagen
	if imagen != nil {
		imagenURL, err := s.guardarImagen(imagen)
		if err != nil {
			return nil, fmt.Errorf("error al guardar imagen: %v", err)
		}
		usuario.Imagen = imagenURL
	} else {
		return nil, errors.New("imagen vacía")
	}

	if err := database.DB.Save(&usuario).Error; err != nil {
		return nil, err
	}

	return &usuario, nil
}

func (s *UsuarioService) guardarImagen(fileHeader *multipart.FileHeader) (string, error) {
	// Crear directorio si no existe
	uploadDir := "uploads/usuarios"
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

	return fmt.Sprintf("%s/usuarios/%s", baseURL, filename), nil
}
