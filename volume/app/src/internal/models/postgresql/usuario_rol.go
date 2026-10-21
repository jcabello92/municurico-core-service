package postgresql

import (
	"github.com/google/uuid"
)

type UsuarioRol struct {
	Base

	UsuarioID uuid.UUID `gorm:"type:uuid;not null;index" json:"usuario_id"`
	RolID     uuid.UUID `gorm:"type:uuid;not null;index" json:"rol_id"`

	// Relaciones GORM
	Usuario *Usuario `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"usuario,omitempty"`
	Rol     *Rol     `gorm:"foreignKey:RolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rol,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (UsuarioRol) TableName() string {
	return "usuario_roles"
}
