package shared

import (
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/Bhinneka/user-service/helper"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTExtractor(t *testing.T) {
	tokenStr := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6dHJ1ZSwiZGlkIjoiYzBiNGQxYjRjNDQ3NCIsImRsaSI6IldFQiIsImVtYWlsIjoiYmVuam9AbWFpbGluYXRvci5jb20iLCJpYXQiOjE1NDE2NTUwMTgsImlzcyI6ImJoaW5uZWthLmNvbSIsInN0YWZmIjpmYWxzZSwic3ViIjoiVVNSMTgxMTAxIn0.t89BKgqikjLxAAxy7RzbYqR_TAgwcWNbrb2a7Hyqth170J1bvZ1BODIBkYB8DQHr1_o5vYK7tjQSPVZ0kzWjK-pG3YtyaTPXeUQ4SeHvup5vlG4DI6iXYDQuupS-UkyfDcePPYaUFwOZggR35Pe0bL-INIEDaV0v_6hzTJ9FGGSSwT2Q6zxRRr3BRhktlrVbL-C-1Axk7xkTPiKiaJCyDc1gx9NeC8zHBPB9VIuLDK4tpJXYatvuXJcYonudGdK-LcxMyxgGaYanY6rzQpSmVAYNNJbU4lrYFTzWuIZh3nfIFJW5kL2tifttdm9wblb5dhm56GUhwEGkQ89O7LuR9w"
	test := []struct {
		name      string
		token     string
		wantError bool
	}{
		{
			name:      helper.TextTestCase1,
			token:     tokenStr,
			wantError: false,
		},
		{
			name:      helper.TextTestCase2,
			token:     "",
			wantError: true,
		},
		{
			name:      helper.TextTestCase3,
			token:     `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6ImJoaW5uZWthLW1pY3Jvc2VydmljZXMtYjEzNzE0LTUzMTIxMTUiLCJhdXRob3Jpc2VkIjp0cnVlLCJkaWQiOiJjMGI0ZDFiNGM0NDc0IiwiZGxpIjoiV0VCIiwiaWF0IjoxNTQ0NTQyOTYwLCJpc3MiOiJiaGlubmVrYS5jb20iLCJzdWIiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1In0.IgXWVme1braEjXuGpJ-faz6UpTndH24k95TIkI_kj6RNEGQzyshByHSn377tzY3-SkA6MMbo5FIl8U8l4JP3q1oCY2n_2jWxQM9wzO-TlUhZJKoOCvNTlYzuzqYHnNz9GXiATfB4zqF_HHHdrHMQiVUYiUJVQLhjcxtgqrLLxUo`,
			wantError: true,
		},
		{
			name:      helper.TextTestCase4,
			token:     "asd" + tokenStr,
			wantError: true,
		},

		{
			name:      helper.TextTestCase5,
			token:     `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJiaGlubmVrYS1taWNyb3NlcnZpY2VzLWIxMzcxNC01MzEyMTE1IiwiYXV0aG9yaXNlZCI6dHJ1ZSwiZGlkIjoiYzBiNGQxYjRjNDQ3NCIsImRsaSI6IldFQiIsImVtYWlsIjoidGVzdGRldjg4MEBnbWFpbC5jb20iLCJleHAiOjE1OTU1ODY0MTksImlhdCI6MTU5NTU3OTIxOSwiaXNzIjoiYmhpbm5la2EuY29tIiwianRpIjoiNTU4NzI0NzAxMjExMGY5MTRkMGE5NjQ2YWVmN2FjOTE1NDEzZGJiNyIsIm1lbWJlclR5cGUiOiJwZXJzb25hbCIsInNpZ25VcEZyb20iOiIiLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIwMDY5MTAwMTY5NSJ9.F8rRBeAnPkqpkl-aoCO_ys3qbFyBLgFF5S3gyGuA-0-ih55W9Wk2m-kY2-Be7SwEbUoaTJcj3l3DeLUxTvKj0bXEtsUdoTTogLbN8AAi-YaDrUgoFTrJo9bhPdkbsLsViAyb3xJwwB2bciswzKYrMBYj7Q1rB30BM0ev0XGnG8KX3KkFEIMXk7I1ks9NMxpz7vdFGrnm7-BBZCbmGZ-tVFaTKzVuOJMRIw8WAMKCZ4jGPqTUPgKac3k7TMlpZg1FpOuX9jks3eo_MdrKAWu192Fv4ktGsG4fApA0aZFgW3YupwPYwQg48LQleu5_c1CKq78BJk4D7FAfTjaC_FcaCw`,
			wantError: true,
		},
	}

	for _, tt := range test {
		verifyKey, _ := getPublicKey()
		claims, err := JWTExtract(verifyKey, tt.token)
		if !tt.wantError {
			assert.NoError(t, err)
			assert.Equal(t, "USR181101", claims.Subject)
			assert.Equal(t, "c0b4d1b4c4474", claims.DeviceID)
		} else {
			assert.Error(t, err)
		}
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
