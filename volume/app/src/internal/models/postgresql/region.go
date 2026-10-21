package postgresql

type Region struct {
	Base

	Nombre      string `gorm:"type:varchar(100);not null" json:"nombre"`
	Abreviatura string `gorm:"type:varchar(10);not null" json:"abreviatura"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Region) TableName() string {
	return "regiones"
}
