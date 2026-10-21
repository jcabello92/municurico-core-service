package enums

import "errors"

type EstadoSesion string

const (
    EstadoSesionActiva     EstadoSesion = "ACTIVA"
    EstadoSesionExpirada   EstadoSesion = "EXPIRADA"
    EstadoSesionFinalizada EstadoSesion = "FINALIZADA"
    EstadoSesionRevocada   EstadoSesion = "REVOCADA"
)

// String retorna el valor del enum como string.
func (e EstadoSesion) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e EstadoSesion) EsValido() bool {
    switch e {
    case EstadoSesionActiva, EstadoSesionExpirada, EstadoSesionFinalizada, EstadoSesionRevocada:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e EstadoSesion) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para estado sesión")
	}
	return nil
}