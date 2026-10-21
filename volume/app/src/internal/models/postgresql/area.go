package postgresql

import (
	"github.com/google/uuid"

	"core-service/internal/enums"
)

type Area struct {
	Base

	AreaPadreID             *uuid.UUID                    `gorm:"type:uuid;index" json:"area_padre_id,omitempty"`
	Nombre                  string                        `gorm:"type:varchar(100);not null" json:"nombre"`
	CargoID                 uuid.UUID                     `gorm:"type:uuid;not null;index" json:"cargo_id"`
	ResponsableTitularID    uuid.UUID                     `gorm:"type:uuid;not null;index" json:"responsable_titular_id"`
	ResponsableSuplenteID   *uuid.UUID                    `gorm:"type:uuid;index" json:"responsable_suplente_id,omitempty"`
	ResponsableSeleccionado enums.ResponsableSeleccionado `gorm:"type:varchar(10);not null" json:"responsable_seleccionado"`

	// Relaciones GORM
	AreaPadre          *Area    `gorm:"foreignKey:AreaPadreID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"area_padre,omitempty"`
	Cargo              *Cargo   `gorm:"foreignKey:CargoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"cargo,omitempty"`
	ResponsableTitular *Miembro `gorm:"foreignKey:ResponsableTitularID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"responsable_titular,omitempty"`
	ResponsableSuplente *Miembro `gorm:"foreignKey:ResponsableSuplenteID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"responsable_suplente,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Area) TableName() string {
	return "areas"
}
