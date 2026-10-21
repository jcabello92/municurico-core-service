package postgresql

import (
	"time"

	"github.com/google/uuid"

	"core-service/internal/enums"
)

type Sesion struct {
	PrimaryKey

	UsuarioID         uuid.UUID          `gorm:"type:uuid;not null;index" json:"usuario_id"`
	UsuarioIP         string             `gorm:"type:inet;not null" json:"usuario_ip"`
	UserAgent         string             `gorm:"type:text;not null" json:"user_agent"`
	FechaInicio       time.Time          `gorm:"not null" json:"fecha_inicio"`
	FechaExpiracion   time.Time          `gorm:"not null" json:"fecha_expiracion"`
	FechaFinalizacion *time.Time         `json:"fecha_finalizacion,omitempty"`
	Estado            enums.EstadoSesion `gorm:"type:varchar(10);not null" json:"estado"`
	CreatedAt         time.Time          `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time          `gorm:"not null" json:"updated_at"`

	// Relaciones GORM
	Usuario *Usuario `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"usuario,omitempty"`
}

// TableName sobreescribe el nombre de la tabla para adaptarse al nombre real que esta tiene en la base de datos.
func (Sesion) TableName() string {
	return "sesiones"
}
