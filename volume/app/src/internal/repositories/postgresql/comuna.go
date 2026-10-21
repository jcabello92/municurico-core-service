package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	models "core-service/internal/models/postgresql"
)

// ComunaRepository gestiona la persistencia y consultas específicas para la entidad Comuna.
type ComunaRepository struct {
	*BaseRepository[models.Comuna]
}

// NewComunaRepository inicializa y retorna una nueva instancia de ComunaRepository.
func NewComunaRepository(db *gorm.DB) *ComunaRepository {
	return &ComunaRepository{
		BaseRepository: NewBaseRepository[models.Comuna](db),
	}
}

// ListByProvinciaID obtiene todas las comunas pertenecientes a una provincia específica.
func (r *ComunaRepository) ListByProvinciaID(ctx context.Context, provinciaID string) ([]models.Comuna, error) {
	if strings.TrimSpace(provinciaID) == "" {
		return nil, errors.New("el provincia_id no puede estar vacío")
	}

	provinciaUUID, err := uuid.Parse(provinciaID)
	if err != nil {
		return nil, fmt.Errorf("el provincia_id '%s' no es un UUID válido: %w", provinciaID, err)
	}

	var comunas []models.Comuna
	if err := r.db.WithContext(ctx).Where("provincia_id = ?", provinciaUUID).Find(&comunas).Error; err != nil {
		return nil, fmt.Errorf("error al listar comunas por provincia_id %s: %w", provinciaID, err)
	}

	return comunas, nil
}
