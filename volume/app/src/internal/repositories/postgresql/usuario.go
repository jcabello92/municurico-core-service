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

// UsuarioRepository gestiona la persistencia y consultas específicas para la entidad Usuario.
type UsuarioRepository struct {
	*BaseRepository[models.Usuario]
}

// NewUsuarioRepository inicializa y retorna una nueva instancia de UsuarioRepository.
func NewUsuarioRepository(db *gorm.DB) *UsuarioRepository {
	return &UsuarioRepository{
		BaseRepository: NewBaseRepository[models.Usuario](db),
	}
}

// ReadByPersonaID busca un usuario a partir de su persona_id referencial.
func (r *UsuarioRepository) ReadByPersonaID(ctx context.Context, personaID string) (*models.Usuario, error) {
	if strings.TrimSpace(personaID) == "" {
		return nil, errors.New("el persona_id no puede estar vacío")
	}

	personaUUID, err := uuid.Parse(personaID)
	if err != nil {
		return nil, fmt.Errorf("el persona_id '%s' no es un UUID válido: %w", personaID, err)
	}

	var usuario models.Usuario
	if err := r.db.WithContext(ctx).Where("persona_id = ?", personaUUID).First(&usuario).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no se encontró ningún usuario con persona_id %s: %w", personaID, err)
		}
		return nil, fmt.Errorf("error al buscar usuario por persona_id %s: %w", personaID, err)
	}

	return &usuario, nil
}

// UpdateEstado actualiza el estado de un usuario validando el nuevo valor del enum.
func (r *UsuarioRepository) UpdateEstado(ctx context.Context, id string, nuevoEstado enums.EstadoUsuario) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("el identificador (id) no puede estar vacío")
	}

	if err := nuevoEstado.Validar(); err != nil {
		return fmt.Errorf("estado inválido: %w", err)
	}

	usuario, err := r.Read(ctx, id)
	if err != nil {
		return err
	}

	usuario.Estado = nuevoEstado
	_, err = r.Update(ctx, id, usuario)
	return err
}

// UpdateContrasena actualiza el hash de contraseña de un usuario mediante BaseRepository registrando la auditoría.
func (r *UsuarioRepository) UpdateContrasena(ctx context.Context, id string, hashContrasena string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("el identificador (id) no puede estar vacío")
	}
	if strings.TrimSpace(hashContrasena) == "" {
		return errors.New("el hash de contraseña no puede estar vacío")
	}

	usuario, err := r.Read(ctx, id)
	if err != nil {
		return err
	}

	usuario.Contrasena = &hashContrasena
	_, err = r.Update(ctx, id, usuario)
	return err
}
