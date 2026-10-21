package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	models "core-service/internal/models/postgresql"
)

// RegionRepository gestiona la persistencia y consultas específicas para la entidad Region.
type RegionRepository struct {
	*BaseRepository[models.Region]
}

// NewRegionRepository inicializa y retorna una nueva instancia de RegionRepository.
func NewRegionRepository(db *gorm.DB) *RegionRepository {
	return &RegionRepository{
		BaseRepository: NewBaseRepository[models.Region](db),
	}
}

// ReadByAbreviatura busca una región a partir de su abreviatura oficial.
func (r *RegionRepository) ReadByAbreviatura(ctx context.Context, abreviatura string) (*models.Region, error) {
	if strings.TrimSpace(abreviatura) == "" {
		return nil, errors.New("la abreviatura no puede estar vacía")
	}

	var region models.Region
	if err := r.db.WithContext(ctx).Where("abreviatura = ?", strings.TrimSpace(abreviatura)).First(&region).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ninguna región con abreviatura '%s'", abreviatura)
		}
		return nil, fmt.Errorf("error al buscar región por abreviatura: %w", err)
	}

	return &region, nil
}
