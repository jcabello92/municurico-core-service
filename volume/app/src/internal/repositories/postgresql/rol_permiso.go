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

// RolPermisoRepository gestiona la persistencia y consultas específicas para la relación RolPermiso.
type RolPermisoRepository struct {
	*BaseRepository[models.RolPermiso]
}

// NewRolPermisoRepository inicializa y retorna una nueva instancia de RolPermisoRepository.
func NewRolPermisoRepository(db *gorm.DB) *RolPermisoRepository {
	return &RolPermisoRepository{
		BaseRepository: NewBaseRepository[models.RolPermiso](db),
	}
}

// ListByRolID obtiene todas las asignaciones de permisos para un rol específico.
func (r *RolPermisoRepository) ListByRolID(ctx context.Context, rolID string) ([]models.RolPermiso, error) {
	if strings.TrimSpace(rolID) == "" {
		return nil, errors.New("el rol_id no puede estar vacío")
	}

	rolUUID, err := uuid.Parse(rolID)
	if err != nil {
		return nil, fmt.Errorf("el rol_id '%s' no es un UUID válido: %w", rolID, err)
	}

	var items []models.RolPermiso
	if err := r.db.WithContext(ctx).Where("rol_id = ?", rolUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar permisos por rol_id %s: %w", rolID, err)
	}

	return items, nil
}

// ListByPermisoID obtiene todas las asignaciones que incluyen un permiso específico.
func (r *RolPermisoRepository) ListByPermisoID(ctx context.Context, permisoID string) ([]models.RolPermiso, error) {
	if strings.TrimSpace(permisoID) == "" {
		return nil, errors.New("el permiso_id no puede estar vacío")
	}

	permisoUUID, err := uuid.Parse(permisoID)
	if err != nil {
		return nil, fmt.Errorf("el permiso_id '%s' no es un UUID válido: %w", permisoID, err)
	}

	var items []models.RolPermiso
	if err := r.db.WithContext(ctx).Where("permiso_id = ?", permisoUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar asignaciones por permiso_id %s: %w", permisoID, err)
	}

	return items, nil
}

// DeleteByRolIDAndPermisoID elimina la asignación entre un rol y un permiso.
func (r *RolPermisoRepository) DeleteByRolIDAndPermisoID(ctx context.Context, rolID, permisoID string) error {
	if strings.TrimSpace(rolID) == "" || strings.TrimSpace(permisoID) == "" {
		return errors.New("los identificadores de rol y permiso son obligatorios")
	}

	rolUUID, err := uuid.Parse(rolID)
	if err != nil {
		return fmt.Errorf("el rol_id '%s' no es un UUID válido: %w", rolID, err)
	}

	permisoUUID, err := uuid.Parse(permisoID)
	if err != nil {
		return fmt.Errorf("el permiso_id '%s' no es un UUID válido: %w", permisoID, err)
	}

	var item models.RolPermiso
	if err := r.db.WithContext(ctx).Where("rol_id = ? AND permiso_id = ?", rolUUID, permisoUUID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no se encontró asignación para rol %s y permiso %s", rolID, permisoID)
		}
		return fmt.Errorf("error al buscar asignación de permiso: %w", err)
	}

	_, err = r.SoftDelete(ctx, item.ID.String())
	return err
}
