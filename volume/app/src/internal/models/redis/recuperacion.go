package redis

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RecuperacionKeyPrefix es el prefijo base utilizado para las claves de recuperación de contraseña en Redis.
const RecuperacionKeyPrefix = "recuperacion"

// BuildRecuperacionKey genera la clave estandarizada de Redis con el formato recuperacion:token.
func BuildRecuperacionKey(token string) string {
	return fmt.Sprintf("%s:%s", RecuperacionKeyPrefix, token)
}

// Recuperacion representa los datos almacenados en Redis para la recuperación de contraseña.
type Recuperacion struct {
	UsuarioID uuid.UUID `json:"usuario_id"`
	CreatedAt time.Time `json:"created_at"`
}

// BuildKey retorna la clave de Redis para el token de recuperación provisto.
func (r *Recuperacion) BuildKey(token string) string {
	return BuildRecuperacionKey(token)
}
