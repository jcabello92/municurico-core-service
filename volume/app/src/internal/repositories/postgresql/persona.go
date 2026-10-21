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

// PersonaRepository gestiona la persistencia y consultas específicas para la entidad Persona.
type PersonaRepository struct {
	*BaseRepository[models.Persona]
}

// NewPersonaRepository inicializa y retorna una nueva instancia de PersonaRepository.
func NewPersonaRepository(db *gorm.DB) *PersonaRepository {
	return &PersonaRepository{
		BaseRepository: NewBaseRepository[models.Persona](db),
	}
}

// ReadByRunBindex busca una persona a partir del blind index de su RUN.
func (r *PersonaRepository) ReadByRunBindex(ctx context.Context, runBindex string) (*models.Persona, error) {
	if strings.TrimSpace(runBindex) == "" {
		return nil, errors.New("el run_bindex no puede estar vacío")
	}

	var persona models.Persona
	if err := r.db.WithContext(ctx).Where("run_bindex = ?", runBindex).First(&persona).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ninguna persona con run_bindex %s: %w", runBindex, err)
		}
		return nil, fmt.Errorf("error al buscar persona por run_bindex %s: %w", runBindex, err)
	}

	return &persona, nil
}

// ReadByPasaporteBindex busca una persona a partir del blind index de su pasaporte.
func (r *PersonaRepository) ReadByPasaporteBindex(ctx context.Context, pasaporteBindex string) (*models.Persona, error) {
	if strings.TrimSpace(pasaporteBindex) == "" {
		return nil, errors.New("el pasaporte_bindex no puede estar vacío")
	}

	var persona models.Persona
	if err := r.db.WithContext(ctx).Where("pasaporte_bindex = ?", pasaporteBindex).First(&persona).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ninguna persona con pasaporte_bindex %s: %w", pasaporteBindex, err)
		}
		return nil, fmt.Errorf("error al buscar persona por pasaporte_bindex %s: %w", pasaporteBindex, err)
	}

	return &persona, nil
}

// ReadByCorreoElectronicoBindex busca una persona a partir del blind index de su correo electrónico.
func (r *PersonaRepository) ReadByCorreoElectronicoBindex(ctx context.Context, correoBindex string) (*models.Persona, error) {
	if strings.TrimSpace(correoBindex) == "" {
		return nil, errors.New("el correo_electronico_bindex no puede estar vacío")
	}

	var persona models.Persona
	if err := r.db.WithContext(ctx).Where("correo_electronico_bindex = ?", correoBindex).First(&persona).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ninguna persona con correo_electronico_bindex %s: %w", correoBindex, err)
		}
		return nil, fmt.Errorf("error al buscar persona por correo_electronico_bindex %s: %w", correoBindex, err)
	}

	return &persona, nil
}

// ListByComunaID obtiene la colección de personas asociadas a una comuna específica.
func (r *PersonaRepository) ListByComunaID(ctx context.Context, comunaID string) ([]models.Persona, error) {
	if strings.TrimSpace(comunaID) == "" {
		return nil, errors.New("el comuna_id no puede estar vacío")
	}

	comunaUUID, err := uuid.Parse(comunaID)
	if err != nil {
		return nil, fmt.Errorf("el comuna_id '%s' no es un UUID válido: %w", comunaID, err)
	}

	var personas []models.Persona
	if err := r.db.WithContext(ctx).Where("comuna_id = ?", comunaUUID).Find(&personas).Error; err != nil {
		return nil, fmt.Errorf("error al listar personas por comuna_id %s: %w", comunaID, err)
	}

	return personas, nil
}
