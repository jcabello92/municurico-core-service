package postgresql

import (
	"github.com/google/uuid"

	"core-service/internal/enums"
)

type Miembro struct {
	Base

	PersonaID               uuid.UUID           `gorm:"type:uuid;not null;index" json:"persona_id"`
	TelefonoFijo            *string             `gorm:"type:text" json:"telefono_fijo,omitempty"`
	TelefonoMovil           *string             `gorm:"type:text" json:"telefono_movil,omitempty"`
	CorreoElectronico       *string             `gorm:"type:text" json:"correo_electronico,omitempty"`
	CorreoElectronicoBindex *string             `gorm:"type:varchar(64);index" json:"correo_electronico_bindex,omitempty"`
	Estado                  enums.EstadoMiembro `gorm:"type:varchar(10);not null" json:"estado"`

	// Relaciones GORM
	Persona *Persona `gorm:"foreignKey:PersonaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"persona,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Miembro) TableName() string {
	return "miembros"
}
