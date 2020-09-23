package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/pkg/errors"
)

// todo: make an interface and support other types of encryption.
type Encryptor struct {
	settings *Settings
}

func NewEncryptor(settings *Settings) *Encryptor {
	return &Encryptor{settings: settings}
}

func (e *Encryptor) Encrypt(data []byte) ([]byte, error) {
	if !e.settings.TurnedOn {
		return data, nil

	}

	blockCipher, err := e.newAesCipher()
	if err != nil {
		return nil, errors.WithMessage(err, "new AES cipher")
	}

	gcm, err := e.newGCM(blockCipher)
	if err != nil {
		return nil, errors.WithMessage(err, "new GCM")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, errors.WithMessage(err, "read random nonce")
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {
	if !e.settings.TurnedOn {
		return data, nil
	}

	blockCipher, err := e.newAesCipher()
	if err != nil {
		return nil, errors.WithMessage(err, "new aes cipher")
	}

	gcm, err := e.newGCM(blockCipher)
	if err != nil {
		return nil, errors.WithMessage(err, "new GCM")
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.WithMessage(err, "open GCM")
	}
	return plaintext, nil
}

func (e *Encryptor) newAesCipher() (cipher.Block, error) {
	return aes.NewCipher(e.saltBytes())
}

func (e *Encryptor) newGCM(blockCipher cipher.Block) (cipher.AEAD, error) {
	return cipher.NewGCM(blockCipher)
}

func (e *Encryptor) saltBytes() []byte {
	return []byte(e.settings.Salt)
}
