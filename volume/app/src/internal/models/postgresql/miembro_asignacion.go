package postgresql

import (
	"github.com/google/uuid"
)

type MiembroAsignacion struct {
	Base

	MiembroID uuid.UUID `gorm:"type:uuid;not null;index" json:"miembro_id"`
	AreaID    uuid.UUID `gorm:"type:uuid;not null;index" json:"area_id"`
	CargoID   uuid.UUID `gorm:"type:uuid;not null;index" json:"cargo_id"`

	// Relaciones GORM
	Miembro *Miembro `gorm:"foreignKey:MiembroID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"miembro,omitempty"`
	Area    *Area    `gorm:"foreignKey:AreaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"area,omitempty"`
	Cargo   *Cargo   `gorm:"foreignKey:CargoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"cargo,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (MiembroAsignacion) TableName() string {
	return "miembro_asignaciones"
}
