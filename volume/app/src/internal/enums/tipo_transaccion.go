package enums

import "errors"

type TipoTransaccion string

const (
    TipoTransaccionInsert  TipoTransaccion = "INSERT"
    TipoTransaccionUpdate  TipoTransaccion = "UPDATE"
    TipoTransaccionDelete  TipoTransaccion = "DELETE"
    TipoTransaccionRestore TipoTransaccion = "RESTORE"
)

// String retorna el valor del enum como string.
func (e TipoTransaccion) String() string {
	return string(e)
}

// EsValido retorna true si el enum es válido.
func (e TipoTransaccion) EsValido() bool {
    switch e {
    case TipoTransaccionInsert, TipoTransaccionUpdate, TipoTransaccionDelete, TipoTransaccionRestore:
        return true
    }
    return false
}

// Validar retorna un error si el enum no es válido.
func (e TipoTransaccion) Validar() error {
	if !e.EsValido() {
		return errors.New("valor no válido para tipo transacción")
	}
	return nil
}