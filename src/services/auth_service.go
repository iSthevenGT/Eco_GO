package services

import (
	"Eco_GO/src/database"
	"Eco_GO/src/dto"
	"Eco_GO/src/models"
	"Eco_GO/src/repositories"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	usuarioRepo *repositories.UsuarioRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		usuarioRepo: repositories.NewUsuarioRepository(),
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (map[string]interface{}, error) {
	var usuario models.Usuario

	// Buscar usuario por Email con todas las relaciones
	if err := database.DB.
		Preload("Telefono").
		Preload("Consumidor").
		Preload("Comerciante").
		Preload("Repartidor").
		Where("Email = ?", req.Email).First(&usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("credenciales inválidas")
		}
		return nil, err
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Contrasena), []byte(req.Contrasena)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token, err := s.GenerateToken(usuario)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"token":     token,
		"rol":       usuario.Rol,
		"idUsuario": usuario.ID,
		"Email":     usuario.Email,
		"nombre":    usuario.Nombre,
		"telefono":  usuario.Telefono,
		"imagen":    usuario.Imagen,
		"message":   "Login exitoso",
	}

	// Agregar información específica del rol
	switch usuario.Rol {
	case "COMERCIANTE":
		if usuario.Comerciante != nil {
			response["nit"] = usuario.Comerciante.NIT
			response["rut"] = usuario.Comerciante.RUT
			response["camara_comercio"] = usuario.Comerciante.CamaraComercio
		}
	case "CONSUMIDOR":
		if usuario.Consumidor != nil {
			response["puntos"] = usuario.Consumidor.Puntos
		}
	case "REPARTIDOR":
		if usuario.Repartidor != nil {
			response["placa"] = usuario.Repartidor.Placa
			response["licencia"] = usuario.Repartidor.Licencia
		}
	}

	return response, nil
}

func (s *AuthService) Register(req dto.RegisterRequest) (*models.Usuario, error) {
	// Verificar si el Email ya existe
	if err := s.usuarioRepo.ValidarEmail(req.Email); err != nil {
		return nil, err
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Contrasena), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Crear usuario base
	usuario := models.Usuario{
		Nombre:     req.Nombre,
		Email:      req.Email,
		Contrasena: string(hashedPassword),
		Rol:        req.Rol,
	}

	usuarioCreado, err := s.usuarioRepo.Crear(&usuario, req.Telefono)
	//                                                               ↑ Quitar punto y coma
	if err != nil {
		return nil, err
	}

	return usuarioCreado, nil
}

func (s *AuthService) GenerateToken(usuario models.Usuario) (string, error) {
	claims := jwt.MapClaims{
		"id":  usuario.ID,
		"rol": usuario.Rol,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("token inválido")
}

func (s *AuthService) GetUserIDFromToken(tokenString string) (uint, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	userIDFloat, ok := (*claims)["id"].(float64)
	if !ok {
		return 0, errors.New("ID de usuario inválido en token")
	}

	return uint(userIDFloat), nil
}

func (s *AuthService) GetRoleFromToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	role, ok := (*claims)["rol"].(string)
	if !ok {
		return "", errors.New("rol inválido en token")
	}

	return role, nil
}

// Función auxiliar para obtener datos completos del usuario
func (s *AuthService) GetUserWithRole(userID uint) (*models.Usuario, error) {
	var usuario models.Usuario

	if err := database.DB.Preload("Telefono").
		Preload("Consumidor").
		Preload("Comerciante").
		Preload("Repartidor").
		First(&usuario, userID).Error; err != nil {
		return nil, err
	}

	return &usuario, nil
}
