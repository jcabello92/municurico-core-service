package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Bindexer genera índices ciegos deterministas mediante HMAC-SHA256 para búsquedas exactas sobre datos cifrados.
type Bindexer struct {
	secret []byte
}

// NewBindexer crea una nueva instancia de Bindexer con el secreto provisto.
func NewBindexer(secret []byte) (*Bindexer, error) {
	if len(secret) == 64 {
		decodedSecret, err := hex.DecodeString(string(secret))
		if err == nil && len(decodedSecret) == 32 {
			secret = decodedSecret
		}
	}

	if len(secret) == 0 {
		return nil, errors.New("la clave secreta para el blind index no puede estar vacía")
	}
	return &Bindexer{secret: secret}, nil
}

// GenerateBindex genera un hash HMAC-SHA256 normalizado en minúsculas y sin espacios de 64 caracteres hexadecimales.
func (b *Bindexer) GenerateBindex(data string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(data))
	if normalized == "" {
		return "", errors.New("el valor para generar el blind index no puede estar vacío")
	}

	h := hmac.New(sha256.New, b.secret)
	if _, err := h.Write([]byte(normalized)); err != nil {
		return "", fmt.Errorf("error al escribir datos en HMAC: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// GenerateBindexRaw genera un hash HMAC-SHA256 del valor exacto sin aplicar normalización de texto.
func (b *Bindexer) GenerateBindexRaw(data string) (string, error) {
	if data == "" {
		return "", errors.New("el valor para generar el blind index no puede estar vacío")
	}

	h := hmac.New(sha256.New, b.secret)
	if _, err := h.Write([]byte(data)); err != nil {
		return "", fmt.Errorf("error al escribir datos en HMAC: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
