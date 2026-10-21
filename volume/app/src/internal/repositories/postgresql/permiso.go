package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"core-service/internal/enums"
	models "core-service/internal/models/postgresql"
)

// PermisoRepository gestiona la persistencia y consultas específicas para la entidad Permiso.
type PermisoRepository struct {
	*BaseRepository[models.Permiso]
}

// NewPermisoRepository inicializa y retorna una nueva instancia de PermisoRepository.
func NewPermisoRepository(db *gorm.DB) *PermisoRepository {
	return &PermisoRepository{
		BaseRepository: NewBaseRepository[models.Permiso](db),
	}
}

// ReadBySlug busca un permiso a partir del slug de la aplicación y el slug del permiso.
func (r *PermisoRepository) ReadBySlug(ctx context.Context, aplicacionSlug enums.AplicacionSlug, slug string) (*models.Permiso, error) {
	if err := aplicacionSlug.Validar(); err != nil {
		return nil, fmt.Errorf("aplicación slug no válida: %w", err)
	}
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("el slug no puede estar vacío")
	}

	var permiso models.Permiso
	if err := r.db.WithContext(ctx).Where("aplicacion_slug = ? AND slug = ?", aplicacionSlug, slug).First(&permiso).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ningún permiso con slug '%s' en la aplicación '%s': %w", slug, aplicacionSlug, err)
		}
		return nil, fmt.Errorf("error al buscar permiso por slug: %w", err)
	}

	return &permiso, nil
}

// ListByAplicacionSlug obtiene todos los permisos pertenecientes a una aplicación determinada.
func (r *PermisoRepository) ListByAplicacionSlug(ctx context.Context, aplicacionSlug enums.AplicacionSlug) ([]models.Permiso, error) {
	if err := aplicacionSlug.Validar(); err != nil {
		return nil, fmt.Errorf("aplicación slug no válida: %w", err)
	}

	var permisos []models.Permiso
	if err := r.db.WithContext(ctx).Where("aplicacion_slug = ?", aplicacionSlug).Find(&permisos).Error; err != nil {
		return nil, fmt.Errorf("error al listar permisos para la aplicación '%s': %w", aplicacionSlug, err)
	}

	return permisos, nil
}
