package token

import (
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/Bhinneka/user-service/helper"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestValidateTokenIgnoreExpiration(t *testing.T) {
	expiredToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0Ijp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiZXhwIjoxNTIxMDIyNTI0LCJpYXQiOjE1MjEwMjI0NjQsImlzcyI6ImJoaW5uZWthLmNvbSIsInN1YiI6IlVTUjE4MDIxODE3NSJ9.gw_tfbPHq6XIWuVT3ksHFovWkYZteUuuLGepkGeAnAGP41pF_AlEhcB120Jao1FOi74li7f3ab6kN_hBcMNMuYmhlkP4pa78QKFCY7uZpLUm6LIc2AOHf1VRm0poQnvH0AnDDw1_bU8NFe0GKr48Cf88934txTCJRQ75Sw4pnbk"
	verifyKey, err := getPublicKey()
	assert.NoError(t, err)
	claims, err := VerifyTokenIgnoreExpiration(verifyKey, expiredToken)
	assert.NoError(t, err)
	assert.Equal(t, "USR180218175", claims.Subject)
	assert.Equal(t, "c0b4d1b4c4474", claims.DeviceID)

	_, err2 := VerifyTokenIgnoreExpiration(verifyKey, "asd")
	assert.Error(t, err2)
}

func TestValidateTokenInvalid(t *testing.T) {
	test := []struct {
		name  string
		token string
	}{

		{
			name:  helper.TextTestCase1,
			token: "seyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0Ijp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiZXhwIjoxNTIxMDIyNTI0LCJpYXQiOjE1MjEwMjI0NjQsImlzcyI6ImJoaW5uZWthLmNvbSIsInN1YiI6IlVTUjE4MDIxODE3NSJ9.gw_tfbPHq6XIWuVT3ksHFovWkYZteUuuLGepkGeAnAGP41pF_AlEhcB120Jao1FOi74li7f3ab6kN_hBcMNMuYmhlkP4pa78QKFCY7uZpLUm6LIc2AOHf1VRm0poQnvH0AnDDw1_bU8NFe0GKr48Cf88934txTCJRQ75Sw4pnbk",
		},
		{
			name:  helper.TextTestCase2,
			token: `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0Ijp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiZXhwIjoxNTIxMDIyNTI0LCJpYXQiOjE1MjEwMjI0NjQsImlzcyI6ImJoaW5uZWthLmNvbSIsInN1YiI6IlVTUjE4MDIxODE3NSJ9`,
		},
	}

	for _, tt := range test {
		verifyKey, err := getPublicKey()
		assert.NoError(t, err)
		_, err = VerifyTokenIgnoreExpiration(verifyKey, tt.token)
		assert.Error(t, err)
	}
}

const PublicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvSd5tzIKV7rlA4uKn1vl
9cg9hNzjWzbT6cvJy5TED2pilUw6LJZhV+ieV08BX2eoG17ygbs8qs7jAcHPzMWw
MGCIayy8XBNG36diPV9ukFdpLczeov0f6gP093w/C2Y6cLQRN3iBlToZKIR6qf0i
PoFMaqiFa8Ys2OmeEdL2egNm+IxGXxyRB9NOwWGjvt5w7PC41+iIGA/AV9EH7FVe
7bcnBsSGXy3kCTneI/X0pcZq1M7cYEPvzXOtq35xzDrmMSoSPo3O06GyPZNA7S4A
iMpw83U1XNmUsVq7lpXP6sROuxEmPfIVunz13DqVZXOTrtkJoONSgNFJ0VbLKUwb
eQIDAQAB
-----END PUBLIC KEY-----
`

func getPublicKey() (*rsa.PublicKey, error) {
	r := strings.NewReader(PublicKey)
	verifyBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}
	return verifyKey, nil
}
