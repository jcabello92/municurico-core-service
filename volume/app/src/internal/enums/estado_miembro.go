package enums

import "errors"

type EstadoMiembro string

const (
    EstadoMiembroActivo      EstadoMiembro = "ACTIVO"
    EstadoMiembroDesactivado EstadoMiembro = "DESACTIVADO"
)

// String retorna el valor del enum como string.
func (e EstadoMiembro) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e EstadoMiembro) EsValido() bool {
    switch e {
    case EstadoMiembroActivo, EstadoMiembroDesactivado:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e EstadoMiembro) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para estado miembro")
	}
	return nil
}