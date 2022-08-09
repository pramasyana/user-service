package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAES(t *testing.T) {

	t.Run("should success test aes", func(t *testing.T) {
		a := NewAES("bhinneka45rty123")

		expectedDecryptedText := "wuriyanto"

		encryptedText, err := a.Encrypt("wuriyanto")

		assert.NoError(t, err, "encrypting text should not err")

		decryptedText, err := a.Decrypt(encryptedText)

		assert.NoError(t, err, "decrypting text should not err")

		assert.Equal(t, expectedDecryptedText, decryptedText, "decryptedText should equal expected")

		decryptedText, err = a.Decrypt("asd")

		assert.Error(t, err, decryptedText)

	})

	t.Run("should fail test aes", func(t *testing.T) {
		a := NewAES("bhinneka45rty12")

		_, err := a.Encrypt("wuriyanto2")

		assert.Error(t, err, "encrypting text should err")
	})

}
