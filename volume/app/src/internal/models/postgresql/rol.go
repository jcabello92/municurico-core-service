package postgresql

import (
	"core-service/internal/enums"
)

type Rol struct {
	Base

	AplicacionSlug enums.AplicacionSlug `gorm:"type:varchar(100);not null;index" json:"aplicacion_slug"`
	Nombre         string               `gorm:"type:varchar(100);not null" json:"nombre"`
	Descripcion    *string              `gorm:"type:text" json:"descripcion,omitempty"`
	Slug           string               `gorm:"type:varchar(100);not null;index" json:"slug"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Rol) TableName() string {
	return "roles"
}
