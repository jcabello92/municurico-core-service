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

// MiembroAsignacionRepository gestiona la persistencia y consultas específicas para la asignación de miembros a áreas y cargos.
type MiembroAsignacionRepository struct {
	*BaseRepository[models.MiembroAsignacion]
}

// NewMiembroAsignacionRepository inicializa y retorna una nueva instancia de MiembroAsignacionRepository.
func NewMiembroAsignacionRepository(db *gorm.DB) *MiembroAsignacionRepository {
	return &MiembroAsignacionRepository{
		BaseRepository: NewBaseRepository[models.MiembroAsignacion](db),
	}
}

// ListByMiembroID obtiene todas las asignaciones activas asociadas a un miembro específico.
func (r *MiembroAsignacionRepository) ListByMiembroID(ctx context.Context, miembroID string) ([]models.MiembroAsignacion, error) {
	if strings.TrimSpace(miembroID) == "" {
		return nil, errors.New("el miembro_id no puede estar vacío")
	}

	miembroUUID, err := uuid.Parse(miembroID)
	if err != nil {
		return nil, fmt.Errorf("el miembro_id '%s' no es un UUID válido: %w", miembroID, err)
	}

	var items []models.MiembroAsignacion
	if err := r.db.WithContext(ctx).Where("miembro_id = ?", miembroUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar asignaciones por miembro_id %s: %w", miembroID, err)
	}

	return items, nil
}

// ListByAreaID obtiene todas las asignaciones activas asociadas a un área específica.
func (r *MiembroAsignacionRepository) ListByAreaID(ctx context.Context, areaID string) ([]models.MiembroAsignacion, error) {
	if strings.TrimSpace(areaID) == "" {
		return nil, errors.New("el area_id no puede estar vacío")
	}

	areaUUID, err := uuid.Parse(areaID)
	if err != nil {
		return nil, fmt.Errorf("el area_id '%s' no es un UUID válido: %w", areaID, err)
	}

	var items []models.MiembroAsignacion
	if err := r.db.WithContext(ctx).Where("area_id = ?", areaUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar asignaciones por area_id %s: %w", areaID, err)
	}

	return items, nil
}

// ListByCargoID obtiene todas las asignaciones activas asociadas a un cargo específico.
func (r *MiembroAsignacionRepository) ListByCargoID(ctx context.Context, cargoID string) ([]models.MiembroAsignacion, error) {
	if strings.TrimSpace(cargoID) == "" {
		return nil, errors.New("el cargo_id no puede estar vacío")
	}

	cargoUUID, err := uuid.Parse(cargoID)
	if err != nil {
		return nil, fmt.Errorf("el cargo_id '%s' no es un UUID válido: %w", cargoID, err)
	}

	var items []models.MiembroAsignacion
	if err := r.db.WithContext(ctx).Where("cargo_id = ?", cargoUUID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("error al listar asignaciones por cargo_id %s: %w", cargoID, err)
	}

	return items, nil
}

// DeleteByMiembroAreaCargoID elimina la asignación específica entre un miembro, un área y un cargo.
func (r *MiembroAsignacionRepository) DeleteByMiembroAreaCargoID(ctx context.Context, miembroID, areaID, cargoID string) error {
	if strings.TrimSpace(miembroID) == "" || strings.TrimSpace(areaID) == "" || strings.TrimSpace(cargoID) == "" {
		return errors.New("los identificadores de miembro, área y cargo son obligatorios")
	}

	miembroUUID, err := uuid.Parse(miembroID)
	if err != nil {
		return fmt.Errorf("el miembro_id '%s' no es un UUID válido: %w", miembroID, err)
	}

	areaUUID, err := uuid.Parse(areaID)
	if err != nil {
		return fmt.Errorf("el area_id '%s' no es un UUID válido: %w", areaID, err)
	}

	cargoUUID, err := uuid.Parse(cargoID)
	if err != nil {
		return fmt.Errorf("el cargo_id '%s' no es un UUID válido: %w", cargoID, err)
	}

	var item models.MiembroAsignacion
	if err := r.db.WithContext(ctx).Where("miembro_id = ? AND area_id = ? AND cargo_id = ?", miembroUUID, areaUUID, cargoUUID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no se encontró asignación para miembro %s, área %s y cargo %s", miembroID, areaID, cargoID)
		}
		return fmt.Errorf("error al buscar asignación de miembro: %w", err)
	}

	_, err = r.SoftDelete(ctx, item.ID.String())
	return err
}
