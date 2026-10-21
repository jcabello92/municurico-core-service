package enums

import "errors"

type EstadoUsuario string

const (
    EstadoUsuarioPendiente   EstadoUsuario = "PENDIENTE"
    EstadoUsuarioActivo      EstadoUsuario = "ACTIVO"
    EstadoUsuarioBloqueado   EstadoUsuario = "BLOQUEADO"
    EstadoUsuarioSuspendido  EstadoUsuario = "SUSPENDIDO"
    EstadoUsuarioDesactivado EstadoUsuario = "DESACTIVADO"
)

// String retorna el valor del enum como string.
func (e EstadoUsuario) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e EstadoUsuario) EsValido() bool {
    switch e {
    case EstadoUsuarioPendiente, EstadoUsuarioActivo, EstadoUsuarioBloqueado, EstadoUsuarioSuspendido, EstadoUsuarioDesactivado:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e EstadoUsuario) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para estado usuario")
	}
	return nil
}