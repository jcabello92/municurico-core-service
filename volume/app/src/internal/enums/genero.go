package enums

import "errors"

type Genero string

const (
    GeneroMasculino Genero = "MASCULINO"
    GeneroFemenino  Genero = "FEMENINO"
    GeneroOtro      Genero = "OTRO"
)

// String retorna el valor del enum como string.
func (e Genero) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e Genero) EsValido() bool {
    switch e {
    case GeneroMasculino, GeneroFemenino, GeneroOtro:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e Genero) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para género")
	}
	return nil
}