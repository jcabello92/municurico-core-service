package jsonb

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Modificacion representa los valores anterior y actual para un campo auditado.
type Modificacion struct {
	Anterior any `json:"anterior"`
	Actual   any `json:"actual"`
}

// Modificaciones representa el conjunto de modificaciones de un registro donde la clave es el nombre del campo.
type Modificaciones map[string]Modificacion

// Value convierte Modificaciones a un formato compatible con la base de datos (driver.Value).
func (m Modificaciones) Value() (driver.Value, error) {
	if m == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(m)
}

// Scan analiza el valor JSONB recuperado de la base de datos en la estructura Modificaciones.
func (m *Modificaciones) Scan(value any) error {
	if value == nil {
		*m = make(Modificaciones)
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("tipo no compatible para Modificaciones")
	}

	*m = make(Modificaciones)
	return json.Unmarshal(bytes, m)
}
