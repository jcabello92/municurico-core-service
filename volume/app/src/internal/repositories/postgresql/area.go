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

// AreaRepository gestiona la persistencia y consultas específicas para la entidad Area.
type AreaRepository struct {
	*BaseRepository[models.Area]
}

// NewAreaRepository inicializa y retorna una nueva instancia de AreaRepository.
func NewAreaRepository(db *gorm.DB) *AreaRepository {
	return &AreaRepository{
		BaseRepository: NewBaseRepository[models.Area](db),
	}
}

// ListByAreaPadreID obtiene todas las subáreas directas que dependen de un área padre.
func (r *AreaRepository) ListByAreaPadreID(ctx context.Context, areaPadreID string) ([]models.Area, error) {
	if strings.TrimSpace(areaPadreID) == "" {
		return nil, errors.New("el area_padre_id no puede estar vacío")
	}

	padreUUID, err := uuid.Parse(areaPadreID)
	if err != nil {
		return nil, fmt.Errorf("el area_padre_id '%s' no es un UUID válido: %w", areaPadreID, err)
	}

	var areas []models.Area
	if err := r.db.WithContext(ctx).Where("area_padre_id = ?", padreUUID).Find(&areas).Error; err != nil {
		return nil, fmt.Errorf("error al listar subáreas por area_padre_id %s: %w", areaPadreID, err)
	}

	return areas, nil
}

// ListAreasRaiz obtiene las áreas de nivel superior que no poseen área padre asignada.
func (r *AreaRepository) ListAreasRaiz(ctx context.Context) ([]models.Area, error) {
	var areas []models.Area
	if err := r.db.WithContext(ctx).Where("area_padre_id IS NULL").Find(&areas).Error; err != nil {
		return nil, fmt.Errorf("error al listar áreas raíz: %w", err)
	}
	return areas, nil
}
