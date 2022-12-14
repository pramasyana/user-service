package token

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/Bhinneka/user-service/src/auth/v1/repo"
	"github.com/Bhinneka/user-service/src/shared/mocks"
	jwt "github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {
	claims := Claim{"bhinneka", "PWT", "M1", "ts615278cc", "WEB", true, true, true, "wuriyanto@gmail.com", "personal", "sturgeon", "token here"}

	tokenAge, _ := time.ParseDuration("1m")
	refreshTokenAge, _ := time.ParseDuration("1m")

	spTokenAge, _ := time.ParseDuration("2m")
	spRefreshTokenAge, _ := time.ParseDuration("4m")

	fredis := mocks.InitFakeRedis()
	loginSessionRedisRepo := repo.NewLoginSessionRepositoryRedis(fredis)

	validPrivateKey, _ := getPrivateKey(PrivateKey)
	t.Run("Test Generate", func(t *testing.T) {

		var loginSessionRepo repo.LoginSessionRepository
		gen := NewJwtGenerator(validPrivateKey, tokenAge, refreshTokenAge, spTokenAge, spRefreshTokenAge, loginSessionRepo, "email@bhinneka.com")

		tokenResult := <-gen.GenerateAccessToken(claims)

		assert.NoError(t, tokenResult.Error)
	})

	t.Run("Test Generate Anonymous", func(t *testing.T) {
		gen := NewJwtGenerator(validPrivateKey, tokenAge, refreshTokenAge, spTokenAge, spRefreshTokenAge, loginSessionRedisRepo, "email@bhinneka.com")
		tokenResult := <-gen.GenerateAnonymous(context.Background())
		assert.NoError(t, tokenResult.Error)

		assert.NoError(t, tokenResult.Error)
	})
}

const PrivateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCoqzL5JrMzed4tb8uEoLKd42EOsYmb0HpbicGt/OUeJxaHtt59
Ew0BbpreBeiuugXweEa5xctQOxGYr27h4ZOnR0hWSi+h5Y35CKzMEmZnzQwzQphg
qww0U+e9/OAvVfCW1xWvVFr0WbhIRn+w/9DUvp+6jKz3fIj3yQaHWVMMNQIDAQAB
AoGAKd7d94XI5JVzNxpSjmkKDjHc7TXbcEevqDupTdTC19piOGyIDMqG5v0bCtSy
r3VUdh6ViBZ240LWmm2qe/5wlaUSDxGAhQg78cVP9L157hC15vOOckMjcJyuVCpR
Wew61HP3JNPB3dsk8P/fjPwgXzsH9L0zIoT0Krw85TbY8UECQQDsRoyWMPu5V6Sa
kiQrr2hZ+weCRH9Q6yd97UhxPAgSZZswedn97uF5T3GdKvoLpeyUrpBe7x3Y3Ciz
uN+SGBPlAkEAtr/SncRbSqUbdrc5BcAVnFNwLXJc2bN477z7shs0OIAha0rcq/8Q
jC5M5oh8jIM5bltZ5t8CWHbrfryVSwsyEQJACnkiGDI5pkCNSlC6C7mtvXdUIOEa
Z6LU0E8pS+OmU/JvC5oLIKdrFS6BUb8q8EM9lmWafqrIvukbYMQMHPS2RQJASWYN
35PH3tkliK7aVjbp9xmECpzOMhnlTtSmesh2VuMPiRpOOz58lPDbrhPPglgKLwq9
tv6G4KUSvJpdlABxIQJARw4I/XUNdd2ko+gSkEHjwjKg4LlYNydHHGk6RYc1S85L
U8PhZfO17Ul2d9ROvFHx75slASSgWHEnPUF7gphhUA==
-----END RSA PRIVATE KEY-----
`

func getPrivateKey(privateKey string) (*rsa.PrivateKey, error) {
	r := strings.NewReader(privateKey)
	signBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}
	return signKey, nil
}
