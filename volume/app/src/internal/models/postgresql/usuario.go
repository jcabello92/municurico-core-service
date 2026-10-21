package postgresql

import (
	"github.com/google/uuid"

	"core-service/internal/enums"
)

type Usuario struct {
	Base

	PersonaID  uuid.UUID           `gorm:"type:uuid;not null;index" json:"persona_id"`
	Contrasena *string             `gorm:"type:text" json:"contrasena,omitempty"`
	Estado     enums.EstadoUsuario `gorm:"type:varchar(10);not null" json:"estado"`
	EsSistema  bool                `gorm:"not null" json:"es_sistema"`

	// Relaciones GORM
	Persona *Persona `gorm:"foreignKey:PersonaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"persona,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Usuario) TableName() string {
	return "usuarios"
}
