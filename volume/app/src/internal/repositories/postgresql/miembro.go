package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"core-service/internal/enums"
	models "core-service/internal/models/postgresql"
)

// MiembroRepository gestiona la persistencia y consultas específicas para la entidad Miembro.
type MiembroRepository struct {
	*BaseRepository[models.Miembro]
}

// NewMiembroRepository inicializa y retorna una nueva instancia de MiembroRepository.
func NewMiembroRepository(db *gorm.DB) *MiembroRepository {
	return &MiembroRepository{
		BaseRepository: NewBaseRepository[models.Miembro](db),
	}
}

// ReadByPersonaID busca un miembro a partir del identificador de su persona asociada.
func (r *MiembroRepository) ReadByPersonaID(ctx context.Context, personaID string) (*models.Miembro, error) {
	if strings.TrimSpace(personaID) == "" {
		return nil, errors.New("el persona_id no puede estar vacío")
	}

	personaUUID, err := uuid.Parse(personaID)
	if err != nil {
		return nil, fmt.Errorf("el persona_id '%s' no es un UUID válido: %w", personaID, err)
	}

	var miembro models.Miembro
	if err := r.db.WithContext(ctx).Where("persona_id = ?", personaUUID).First(&miembro).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ningún miembro con persona_id %s: %w", personaID, err)
		}
		return nil, fmt.Errorf("error al buscar miembro por persona_id %s: %w", personaID, err)
	}

	return &miembro, nil
}

// ReadByCorreoElectronicoBindex busca un miembro a partir del blind index de su correo electrónico institucional.
func (r *MiembroRepository) ReadByCorreoElectronicoBindex(ctx context.Context, correoBindex string) (*models.Miembro, error) {
	if strings.TrimSpace(correoBindex) == "" {
		return nil, errors.New("el correo_electronico_bindex no puede estar vacío")
	}

	var miembro models.Miembro
	if err := r.db.WithContext(ctx).Where("correo_electronico_bindex = ?", correoBindex).First(&miembro).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ningún miembro con correo_electronico_bindex %s: %w", correoBindex, err)
		}
		return nil, fmt.Errorf("error al buscar miembro por correo_electronico_bindex: %w", err)
	}

	return &miembro, nil
}

// UpdateEstado actualiza el estado de un miembro validando el enum.
func (r *MiembroRepository) UpdateEstado(ctx context.Context, id string, nuevoEstado enums.EstadoMiembro) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("el identificador (id) no puede estar vacío")
	}

	if err := nuevoEstado.Validar(); err != nil {
		return fmt.Errorf("estado de miembro inválido: %w", err)
	}

	miembro, err := r.Read(ctx, id)
	if err != nil {
		return err
	}

	miembro.Estado = nuevoEstado
	_, err = r.Update(ctx, id, miembro)
	return err
}
