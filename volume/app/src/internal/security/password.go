package security

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher administra el hashing y verificación de contraseñas mediante el algoritmo Bcrypt.
type PasswordHasher struct {
	cost int
}

// NewPasswordHasher inicializa una nueva instancia de PasswordHasher con el costo computacional especificado.
func NewPasswordHasher(cost int) *PasswordHasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = bcrypt.DefaultCost
	}
	return &PasswordHasher{cost: cost}
}

// HashPassword genera un hash seguro de la contraseña incluyendo salt aleatorio automático.
func (p *PasswordHasher) HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("la contraseña a hashear no puede estar vacía")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return "", fmt.Errorf("fallo al generar el hash de la contraseña: %w", err)
	}

	return string(bytes), nil
}

// CheckPassword compara una contraseña en texto plano contra un hash de Bcrypt y retorna true si coinciden.
func (p *PasswordHasher) CheckPassword(hash, password string) bool {
	if hash == "" || password == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
