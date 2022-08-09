package rsa

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/golang-jwt/jwt"
)

const (
	privateKeyPath = "config/rsa/app.rsa"
	publicKeyPath  = "config/rsa/app.rsa.pub"
	appleKeyPath   = "config/rsa/authKey_apple.p8"
)

var (
	// VerifyKey rsa from config
	VerifyKey *rsa.PublicKey

	signKey *rsa.PrivateKey
)

// InitPublicKey return *rsa.PublicKey
func InitPublicKey() (*rsa.PublicKey, error) {
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	VerifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}
	return VerifyKey, nil
}

// InitPrivateKey return *rsa.PrivateKey
func InitPrivateKey() (*rsa.PrivateKey, error) {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}
	return signKey, nil
}

// InitAppleClientSecret return clientSecret
func InitAppleClientSecret(clientID, teamID, keyID string) (string, error) {
	ctx := "InitAppleClientSecret"

	appleBaseURL := golib.GetEnvOrFail(ctx, "find_apple_config_auth_env", "APPLE_AUTH_URL")

	signBytes, err := ioutil.ReadFile(appleKeyPath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(signBytes)
	if block == nil {
		return "", errors.New("empty block after decoding")
	}

	x509Encoded := block.Bytes
	parsedKey, err := x509.ParsePKCS8PrivateKey(x509Encoded)
	if err != nil {
		return "", err
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if !ok {
		return "", errors.New("not ecdsa private key")
	}
	now := time.Now()

	// Create the Claims
	claims := &jwt.StandardClaims{
		Issuer:    teamID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Hour*24*1 - time.Second).Unix(), // 180 days
		Audience:  appleBaseURL,
		Subject:   clientID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["key"] = keyID

	return token.SignedString(ecdsaPrivateKey)
}
