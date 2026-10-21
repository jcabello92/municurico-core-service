package postgresql

import (
	"time"

	"github.com/google/uuid"

	"core-service/internal/enums"
)

type Persona struct {
	Base

	Run                     *string       `gorm:"type:text" json:"run,omitempty"`
	RunBindex               *string       `gorm:"type:varchar(64);index" json:"run_bindex,omitempty"`
	Pasaporte               *string       `gorm:"type:text" json:"pasaporte,omitempty"`
	PasaporteBindex         *string       `gorm:"type:varchar(64);index" json:"pasaporte_bindex,omitempty"`
	Nombre                  string        `gorm:"type:varchar(100);not null" json:"nombre"`
	ApellidoPaterno         *string       `gorm:"type:varchar(100)" json:"apellido_paterno,omitempty"`
	ApellidoMaterno         *string       `gorm:"type:varchar(100)" json:"apellido_materno,omitempty"`
	Genero                  *enums.Genero `gorm:"type:varchar(10)" json:"genero,omitempty"`
	FechaNacimiento         *time.Time    `gorm:"type:date" json:"fecha_nacimiento,omitempty"`
	ComunaID                *uuid.UUID    `gorm:"type:uuid;index" json:"comuna_id,omitempty"`
	Direccion               *string       `gorm:"type:text" json:"direccion,omitempty"`
	TelefonoFijo            *string       `gorm:"type:text" json:"telefono_fijo,omitempty"`
	TelefonoMovil           *string       `gorm:"type:text" json:"telefono_movil,omitempty"`
	CorreoElectronico       *string       `gorm:"type:text" json:"correo_electronico,omitempty"`
	CorreoElectronicoBindex *string       `gorm:"type:varchar(64);index" json:"correo_electronico_bindex,omitempty"`

	// Relaciones GORM
	Comuna *Comuna `gorm:"foreignKey:ComunaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"comuna,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Persona) TableName() string {
	return "personas"
}
