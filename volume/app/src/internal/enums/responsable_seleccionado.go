package enums

import "errors"

type ResponsableSeleccionado string

const (
    ResponsableSeleccionadoTitular  ResponsableSeleccionado = "TITULAR"
    ResponsableSeleccionadoSuplente ResponsableSeleccionado = "SUPLENTE"
)

// String retorna el valor del enum como string.
func (e ResponsableSeleccionado) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e ResponsableSeleccionado) EsValido() bool {
    switch e {
    case ResponsableSeleccionadoTitular, ResponsableSeleccionadoSuplente:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e ResponsableSeleccionado) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para responsable seleccionado")
	}
	return nil
}