package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	models "core-service/internal/models/postgresql"
)

// CargoRepository gestiona la persistencia y consultas específicas para la entidad Cargo.
type CargoRepository struct {
	*BaseRepository[models.Cargo]
}

// NewCargoRepository inicializa y retorna una nueva instancia de CargoRepository.
func NewCargoRepository(db *gorm.DB) *CargoRepository {
	return &CargoRepository{
		BaseRepository: NewBaseRepository[models.Cargo](db),
	}
}

// ReadByNombre busca un cargo a partir de su nombre exacto.
func (r *CargoRepository) ReadByNombre(ctx context.Context, nombre string) (*models.Cargo, error) {
	if strings.TrimSpace(nombre) == "" {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	var cargo models.Cargo
	if err := r.db.WithContext(ctx).Where("nombre = ?", strings.TrimSpace(nombre)).First(&cargo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ningún cargo con nombre '%s'", nombre)
		}
		return nil, fmt.Errorf("error al buscar cargo por nombre: %w", err)
	}

	return &cargo, nil
}
