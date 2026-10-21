package postgresql

import (
	"github.com/google/uuid"
)

type RolPermiso struct {
	Base

	RolID     uuid.UUID `gorm:"type:uuid;not null;index" json:"rol_id"`
	PermisoID uuid.UUID `gorm:"type:uuid;not null;index" json:"permiso_id"`

	// Relaciones GORM
	Rol     *Rol     `gorm:"foreignKey:RolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"rol,omitempty"`
	Permiso *Permiso `gorm:"foreignKey:PermisoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"permiso,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (RolPermiso) TableName() string {
	return "rol_permisos"
}
