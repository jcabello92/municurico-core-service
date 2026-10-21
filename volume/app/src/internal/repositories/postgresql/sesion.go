package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"core-service/internal/enums"
	models "core-service/internal/models/postgresql"
)

// SesionRepository gestiona la persistencia y consultas específicas para la entidad Sesion.
type SesionRepository struct {
	*BaseRepository[models.Sesion]
}

// NewSesionRepository inicializa y retorna una nueva instancia de SesionRepository.
func NewSesionRepository(db *gorm.DB) *SesionRepository {
	return &SesionRepository{
		BaseRepository: NewBaseRepository[models.Sesion](db),
	}
}

// ListByUsuarioID obtiene el historial de sesiones asociadas a un usuario.
func (r *SesionRepository) ListByUsuarioID(ctx context.Context, usuarioID string) ([]models.Sesion, error) {
	if strings.TrimSpace(usuarioID) == "" {
		return nil, errors.New("el usuario_id no puede estar vacío")
	}

	usuarioUUID, err := uuid.Parse(usuarioID)
	if err != nil {
		return nil, fmt.Errorf("el usuario_id '%s' no es un UUID válido: %w", usuarioID, err)
	}

	var sesiones []models.Sesion
	if err := r.db.WithContext(ctx).Where("usuario_id = ?", usuarioUUID).Order("created_at DESC").Find(&sesiones).Error; err != nil {
		return nil, fmt.Errorf("error al listar sesiones para el usuario %s: %w", usuarioID, err)
	}

	return sesiones, nil
}

// UpdateEstado actualiza el estado de una sesión validando el enum.
func (r *SesionRepository) UpdateEstado(ctx context.Context, id string, nuevoEstado enums.EstadoSesion) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("el identificador (id) no puede estar vacío")
	}

	if err := nuevoEstado.Validar(); err != nil {
		return fmt.Errorf("estado de sesión inválido: %w", err)
	}

	sesion, err := r.Read(ctx, id)
	if err != nil {
		return err
	}

	sesion.Estado = nuevoEstado
	if nuevoEstado == enums.EstadoSesionFinalizada || nuevoEstado == enums.EstadoSesionRevocada {
		now := time.Now().UTC()
		sesion.FechaFinalizacion = &now
	}

	_, err = r.Update(ctx, id, sesion)
	return err
}

// FinalizarSesion marca una sesión como finalizada registrando la fecha actual de término.
func (r *SesionRepository) FinalizarSesion(ctx context.Context, id string) error {
	return r.UpdateEstado(ctx, id, enums.EstadoSesionFinalizada)
}
