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

// RolRepository gestiona la persistencia y consultas específicas para la entidad Rol.
type RolRepository struct {
	*BaseRepository[models.Rol]
}

// NewRolRepository inicializa y retorna una nueva instancia de RolRepository.
func NewRolRepository(db *gorm.DB) *RolRepository {
	return &RolRepository{
		BaseRepository: NewBaseRepository[models.Rol](db),
	}
}

// ReadBySlug busca un rol a partir del slug de la aplicación y el slug del rol.
func (r *RolRepository) ReadBySlug(ctx context.Context, aplicacionSlug enums.AplicacionSlug, slug string) (*models.Rol, error) {
	if err := aplicacionSlug.Validar(); err != nil {
		return nil, fmt.Errorf("aplicación slug no válida: %w", err)
	}
	if strings.TrimSpace(slug) == "" {
		return nil, errors.New("el slug no puede estar vacío")
	}

	var rol models.Rol
	if err := r.db.WithContext(ctx).Where("aplicacion_slug = ? AND slug = ?", aplicacionSlug, slug).First(&rol).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ningún rol con slug '%s' en la aplicación '%s': %w", slug, aplicacionSlug, err)
		}
		return nil, fmt.Errorf("error al buscar rol por slug: %w", err)
	}

	return &rol, nil
}

// ListByAplicacionSlug obtiene todos los roles pertenecientes a una aplicación determinada.
func (r *RolRepository) ListByAplicacionSlug(ctx context.Context, aplicacionSlug enums.AplicacionSlug) ([]models.Rol, error) {
	if err := aplicacionSlug.Validar(); err != nil {
		return nil, fmt.Errorf("aplicación slug no válida: %w", err)
	}

	var roles []models.Rol
	if err := r.db.WithContext(ctx).Where("aplicacion_slug = ?", aplicacionSlug).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("error al listar roles para la aplicación '%s': %w", aplicacionSlug, err)
	}

	return roles, nil
}
