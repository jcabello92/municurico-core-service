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

// ProvinciaRepository gestiona la persistencia y consultas específicas para la entidad Provincia.
type ProvinciaRepository struct {
	*BaseRepository[models.Provincia]
}

// NewProvinciaRepository inicializa y retorna una nueva instancia de ProvinciaRepository.
func NewProvinciaRepository(db *gorm.DB) *ProvinciaRepository {
	return &ProvinciaRepository{
		BaseRepository: NewBaseRepository[models.Provincia](db),
	}
}

// ListByRegionID obtiene todas las provincias pertenecientes a una región específica.
func (r *ProvinciaRepository) ListByRegionID(ctx context.Context, regionID string) ([]models.Provincia, error) {
	if strings.TrimSpace(regionID) == "" {
		return nil, errors.New("el region_id no puede estar vacío")
	}

	regionUUID, err := uuid.Parse(regionID)
	if err != nil {
		return nil, fmt.Errorf("el region_id '%s' no es un UUID válido: %w", regionID, err)
	}

	var provincias []models.Provincia
	if err := r.db.WithContext(ctx).Where("region_id = ?", regionUUID).Find(&provincias).Error; err != nil {
		return nil, fmt.Errorf("error al listar provincias por region_id %s: %w", regionID, err)
	}

	return provincias, nil
}
