package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"

	models "core-service/internal/models/redis"
)

// RecuperacionRepository gestiona el almacenamiento y operaciones de tokens de recuperación en Redis.
type RecuperacionRepository struct {
	*BaseRepository[models.Recuperacion]
}

// NewRecuperacionRepository inicializa y retorna una nueva instancia de RecuperacionRepository.
func NewRecuperacionRepository(client *goredis.Client) *RecuperacionRepository {
	return &RecuperacionRepository{
		BaseRepository: NewBaseRepository[models.Recuperacion](client),
	}
}

// Save almacena los datos de recuperación en Redis asignando su tiempo de vida (TTL).
func (r *RecuperacionRepository) Save(ctx context.Context, token string, recuperacion *models.Recuperacion, ttl time.Duration) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("el token de recuperación no puede estar vacío")
	}
	if recuperacion == nil {
		return errors.New("los datos de recuperación no pueden ser nulos")
	}

	key := models.BuildRecuperacionKey(token)
	return r.Set(ctx, key, *recuperacion, ttl)
}

// GetByToken recupera los datos de recuperación asociados al token provisto.
func (r *RecuperacionRepository) GetByToken(ctx context.Context, token string) (*models.Recuperacion, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("el token de recuperación no puede estar vacío")
	}

	key := models.BuildRecuperacionKey(token)
	return r.Get(ctx, key)
}

// DeleteByToken elimina la clave de recuperación correspondiente al token provisto.
func (r *RecuperacionRepository) DeleteByToken(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("el token de recuperación no puede estar vacío")
	}

	key := models.BuildRecuperacionKey(token)
	return r.Delete(ctx, key)
}
