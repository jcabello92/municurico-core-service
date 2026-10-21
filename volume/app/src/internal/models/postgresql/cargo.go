package postgresql

type Cargo struct {
	Base

	NombreMasculino string  `gorm:"type:varchar(100);not null" json:"nombre_masculino"`
	NombreFemenino  *string `gorm:"type:varchar(100)" json:"nombre_femenino,omitempty"`
	NombreNeutro    *string `gorm:"type:varchar(100)" json:"nombre_neutro,omitempty"`
	EsResponsable   bool    `gorm:"not null" json:"es_responsable"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Cargo) TableName() string {
	return "cargos"
}
