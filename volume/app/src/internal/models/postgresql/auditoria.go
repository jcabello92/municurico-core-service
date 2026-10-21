package postgresql

import (
	"time"

	"github.com/google/uuid"

	"core-service/internal/enums"
	"core-service/internal/models/jsonb"
)

type Auditoria struct {
	PrimaryKey

	TipoTransaccion enums.TipoTransaccion `gorm:"type:varchar(10);not null" json:"tipo_transaccion"`
	TablaSlug       string                `gorm:"type:varchar(100);not null" json:"tabla_slug"`
	RegistroID      uuid.UUID             `gorm:"type:uuid;not null;index" json:"registro_id"`
	Modificaciones  jsonb.Modificaciones  `gorm:"type:jsonb;not null" json:"modificaciones"`
	AutorID         uuid.UUID             `gorm:"type:uuid;not null;index" json:"autor_id"`
	AutorIP         *string               `gorm:"type:inet" json:"autor_ip,omitempty"`
	UserAgent       *string               `gorm:"type:text" json:"user_agent,omitempty"`
	CreatedAt       time.Time             `gorm:"not null" json:"created_at"`

	// Relaciones GORM
	Autor *Usuario `gorm:"foreignKey:AutorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"autor,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Auditoria) TableName() string {
	return "auditoria"
}
