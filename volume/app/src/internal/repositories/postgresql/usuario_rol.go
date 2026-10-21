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

// UsuarioRolRepository gestiona la persistencia y consultas específicas para la relación UsuarioRol.
type UsuarioRolRepository struct {
	*BaseRepository[models.UsuarioRol]
}

// NewUsuarioRolRepository inicializa y retorna una nueva instancia de UsuarioRolRepository.
func NewUsuarioRolRepository(db *gorm.DB) *UsuarioRolRepository {
	return &UsuarioRolRepository{
		BaseRepository: NewBaseRepository[models.UsuarioRol](db),
	}
}

// ListByUsuarioID obtiene todas las asignaciones de roles para un usuario específico.
func (r *UsuarioRolRepository) ListByUsuarioID(ctx context.Context, usuarioID string) ([]models.UsuarioRol, error) {
	if strings.TrimSpace(usuarioID) == "" {
		return nil, errors.New("el usuario_id no puede estar vacío")
	}

	usuarioUUID, err := uuid.Parse(usuarioID)
	if err != nil {
		return nil, fmt.Errorf("el usuario_id '%s' no es un UUID válido: %w", usuarioID, err)
	}

	var items []models.UsuarioRol
	if err := r.db.WithContext(ctx).Where("usuario_id = ?", usuarioUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar asignaciones por usuario_id %s: %w", usuarioID, err)
	}

	return items, nil
}

// ListByRolID obtiene todas las asignaciones de un rol específico.
func (r *UsuarioRolRepository) ListByRolID(ctx context.Context, rolID string) ([]models.UsuarioRol, error) {
	if strings.TrimSpace(rolID) == "" {
		return nil, errors.New("el rol_id no puede estar vacío")
	}

	rolUUID, err := uuid.Parse(rolID)
	if err != nil {
		return nil, fmt.Errorf("el rol_id '%s' no es un UUID válido: %w", rolID, err)
	}

	var items []models.UsuarioRol
	if err := r.db.WithContext(ctx).Where("rol_id = ?", rolUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar asignaciones por rol_id %s: %w", rolID, err)
	}

	return items, nil
}

// DeleteByUsuarioIDAndRolID elimina la asignación entre un usuario y un rol.
func (r *UsuarioRolRepository) DeleteByUsuarioIDAndRolID(ctx context.Context, usuarioID, rolID string) error {
	if strings.TrimSpace(usuarioID) == "" || strings.TrimSpace(rolID) == "" {
		return errors.New("los identificadores de usuario y rol son obligatorios")
	}

	usuarioUUID, err := uuid.Parse(usuarioID)
	if err != nil {
		return fmt.Errorf("el usuario_id '%s' no es un UUID válido: %w", usuarioID, err)
	}

	rolUUID, err := uuid.Parse(rolID)
	if err != nil {
		return fmt.Errorf("el rol_id '%s' no es un UUID válido: %w", rolID, err)
	}

	var item models.UsuarioRol
	if err := r.db.WithContext(ctx).Where("usuario_id = ? AND rol_id = ?", usuarioUUID, rolUUID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no se encontró asignación para usuario %s y rol %s", usuarioID, rolID)
		}
		return fmt.Errorf("error al buscar asignación de rol: %w", err)
	}

	_, err = r.SoftDelete(ctx, item.ID.String())
	return err
}
