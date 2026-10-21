package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encryptor proporciona métodos de cifrado y descifrado autenticado utilizando AES-256-GCM.
type Encryptor struct {
	key []byte
}

// NewEncryptor inicializa una nueva instancia de Encryptor verificando que la clave tenga exactamente 32 bytes (256 bits).
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) == 64 {
		decodedKey, err := hex.DecodeString(string(key))
		if err == nil && len(decodedKey) == 32 {
			key = decodedKey
		}
	}

	if len(key) != 32 {
		return nil, fmt.Errorf("la clave para AES-256 debe tener exactamente 32 bytes (recibidos: %d)", len(key))
	}
	return &Encryptor{key: key}, nil
}

// Encrypt cifra un texto plano utilizando AES-256-GCM y retorna el resultado codificado en Base64.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", errors.New("el texto plano a cifrar no puede estar vacío")
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("fallo al inicializar bloque AES: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("fallo al inicializar modo GCM: %w", err)
	}

	// Genera un nonce aleatorio seguro de 12 bytes
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("fallo al generar nonce aleatorio: %w", err)
	}

	// Cifra los datos y añade el tag de autenticación
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt descifra un texto cifrado en Base64 utilizando AES-256-GCM y retorna el texto plano original.
func (e *Encryptor) Decrypt(ciphertextBase64 string) (string, error) {
	if ciphertextBase64 == "" {
		return "", errors.New("el texto cifrado no puede estar vacío")
	}

	data, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", fmt.Errorf("fallo al decodificar texto cifrado desde Base64: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("fallo al inicializar bloque AES: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("fallo al inicializar modo GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("la longitud del texto cifrado es menor que el tamaño requerido del nonce")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("fallo al descifrar o autenticar el texto (datos corruptos o clave incorrecta): %w", err)
	}

	return string(plaintext), nil
}
