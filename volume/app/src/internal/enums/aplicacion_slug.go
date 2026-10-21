package enums

import "errors"

type AplicacionSlug string

const (
    AplicacionSlugWebAcceso   AplicacionSlug = "WEB:ACCESO"
    AplicacionSlugWebAdmin    AplicacionSlug = "WEB:ADMIN"
    AplicacionSlugWebIntranet AplicacionSlug = "WEB:INTRANET"
)

// String retorna el valor del enum como string.
func (e AplicacionSlug) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e AplicacionSlug) EsValido() bool {
    switch e {
    case AplicacionSlugWebAcceso, AplicacionSlugWebAdmin, AplicacionSlugWebIntranet:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e AplicacionSlug) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para aplicación slug")
	}
	return nil
}