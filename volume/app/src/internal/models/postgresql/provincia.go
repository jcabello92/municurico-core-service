package postgresql

import (
	"github.com/google/uuid"
)

type Provincia struct {
	Base

	Nombre   string    `gorm:"type:varchar(100);not null" json:"nombre"`
	RegionID uuid.UUID `gorm:"type:uuid;not null;index" json:"region_id"`

	// Relaciones GORM
	Region *Region `gorm:"foreignKey:RegionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"region,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Provincia) TableName() string {
	return "provincias"
}
