package postgresql

import (
	"github.com/google/uuid"
)

type Comuna struct {
	Base

	Nombre      string    `gorm:"type:varchar(100);not null" json:"nombre"`
	ProvinciaID uuid.UUID `gorm:"type:uuid;not null;index" json:"provincia_id"`

	// Relaciones GORM
	Provincia *Provincia `gorm:"foreignKey:ProvinciaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"provincia,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Comuna) TableName() string {
	return "comunas"
}
