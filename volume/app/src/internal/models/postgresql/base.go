package postgresql

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserIDContextKey es la clave tipada para extraer el ID de usuario desde el context.Context.
type contextKey string

const UserIDContextKey contextKey = "userID"

// PrimaryKey contiene el ID de las tablas.
type PrimaryKey struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
}

// Base contiene la estructura base para las tablas con auditoría.
type Base struct {
	PrimaryKey

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	CreatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy uuid.UUID  `gorm:"type:uuid;not null" json:"updated_by"`
	DeletedBy *uuid.UUID `gorm:"type:uuid" json:"deleted_by,omitempty"`

	// Relaciones GORM
	CreatedByUser *Usuario `gorm:"foreignKey:CreatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"created_by_user,omitempty"`
	UpdatedByUser *Usuario `gorm:"foreignKey:UpdatedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"updated_by_user,omitempty"`
	DeletedByUser *Usuario `gorm:"foreignKey:DeletedBy;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"deleted_by_user,omitempty"`
}

// BeforeCreate genera automáticamente un UUIDv7 antes de insertar el registro en la base de datos.
func (p *PrimaryKey) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID, err = uuid.NewV7()
		if err != nil {
			return err
		}
	}
	return
}

// BeforeCreate para Base ejecuta primero el BeforeCreate de PrimaryKey y luego asigna CreatedBy y UpdatedBy.
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	// 1. Genera el UUID llamando al hook de PrimaryKey
	if err = b.PrimaryKey.BeforeCreate(tx); err != nil {
		return err
	}

	// 2. Extrae el usuario autenticado desde el context.Context
	if tx.Statement.Context != nil {
		if val := tx.Statement.Context.Value(UserIDContextKey); val != nil {
			if currentUserID, ok := val.(uuid.UUID); ok && currentUserID != uuid.Nil {
				if b.CreatedBy == uuid.Nil {
					b.CreatedBy = currentUserID
				}
				if b.UpdatedBy == uuid.Nil {
					b.UpdatedBy = currentUserID
				}
			}
		}
	}

	return
}

// BeforeUpdate actualiza UpdatedBy con el usuario autenticado en cada modificación.
func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Context != nil {
		if val := tx.Statement.Context.Value(UserIDContextKey); val != nil {
			if currentUserID, ok := val.(uuid.UUID); ok && currentUserID != uuid.Nil {
				b.UpdatedBy = currentUserID
				tx.Statement.SetColumn("UpdatedBy", currentUserID)
			}
		}
	}
	return
}

// BeforeDelete persiste DeletedBy en la base de datos antes de que GORM ejecute el soft-delete.
func (b *Base) BeforeDelete(tx *gorm.DB) (err error) {
	if tx.Statement.Context != nil {
		if val := tx.Statement.Context.Value(UserIDContextKey); val != nil {
			if currentUserID, ok := val.(uuid.UUID); ok && currentUserID != uuid.Nil {
				b.DeletedBy = &currentUserID

				// Actualiza la columna deleted_by dentro de la transacción activa
				return tx.Session(&gorm.Session{SkipHooks: true}).
					Model(b).
					Update("deleted_by", currentUserID).Error
			}
		}
	}
	return
}
