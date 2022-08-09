package token

import (
	"context"
	"crypto/rsa"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/Bhinneka/user-service/src/auth/v1/repo"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

// Claim data structure
type Claim struct {
	Issuer      string
	Audience    string
	Subject     string
	DeviceID    string
	DeviceLogin string
	Authorised  bool
	IsAdmin     bool
	IsStaff     bool
	Email       string
	MemberType  string
	SignUpFrom  string
	CustomToken string
}

// AccessToken data structure
type AccessToken struct {
	JTI         string
	AccessToken string
	ExpiredAt   time.Time
}

// AccessTokenResponse data structure
type AccessTokenResponse struct {
	Error       error
	AccessToken AccessToken
}

// jwtGenerator private data structure
type jwtGenerator struct {
	signKey                *rsa.PrivateKey
	tokenAge               time.Duration
	refreshTokenAge        time.Duration
	LoginSessionRepo       repo.LoginSessionRepository
	specialTokenAge        time.Duration
	specialRefreshTokenAge time.Duration
	specialEmail           string
}

// AccessTokenGenerator interface abstraction
type AccessTokenGenerator interface {
	GenerateAccessToken(cl Claim) <-chan AccessTokenResponse
	GenerateAnonymous(ctxReq context.Context) <-chan AccessTokenResponse
}

// NewJwtGenerator function for initializing jwtGenerator object
func NewJwtGenerator(signKey *rsa.PrivateKey, tokenAge, refreshTokenAge, stAge, srtAge time.Duration,
	loginSessionRepo repo.LoginSessionRepository, specialEmail string) AccessTokenGenerator {
	return &jwtGenerator{
		signKey:                signKey,
		tokenAge:               tokenAge,
		refreshTokenAge:        refreshTokenAge,
		LoginSessionRepo:       loginSessionRepo,
		specialTokenAge:        stAge,
		specialRefreshTokenAge: srtAge,
		specialEmail:           specialEmail,
	}
}

// GenerateAccessToken function for generating access token
func (j *jwtGenerator) GenerateAccessToken(cl Claim) <-chan AccessTokenResponse {
	result := make(chan AccessTokenResponse)
	go func() {
		defer close(result)

		now := time.Now().Add(-90 * time.Second)
		var age time.Time

		token := jwt.New(jwt.SigningMethodRS256)
		claims := make(jwt.MapClaims)

		uid := uuid.NewV4()
		mixid := cl.Email + "----" + uid.String()
		jti := helper.GenerateTokenByString(mixid)

		emails := strings.Split(j.specialEmail, ",")
		if golib.StringInSlice(cl.Email, emails) {
			age = now.Add(j.specialTokenAge)
		} else {
			age = now.Add(j.tokenAge)
		}

		claims["jti"] = jti
		claims["iss"] = cl.Issuer
		claims["aud"] = cl.Audience
		claims["exp"] = age.Unix()
		claims["iat"] = now.Unix()
		claims["sub"] = cl.Subject
		claims["did"] = cl.DeviceID    // device id
		claims["dli"] = cl.DeviceLogin // device login
		claims["adm"] = cl.IsAdmin
		claims["staff"] = cl.IsStaff
		claims["authorised"] = cl.Authorised // authorised
		claims["email"] = cl.Email
		claims["memberType"] = cl.MemberType
		claims["signUpFrom"] = cl.SignUpFrom
		claims["customToken"] = cl.CustomToken
		token.Claims = claims

		tokenString, err := token.SignedString(j.signKey)
		if err != nil {
			result <- AccessTokenResponse{Error: err}
			return
		}
		result <- AccessTokenResponse{Error: nil, AccessToken: AccessToken{AccessToken: tokenString, ExpiredAt: age, JTI: jti}}
	}()

	return result
}

// GenerateAnonymous function for default token anonymous data into redis DB
func (j *jwtGenerator) GenerateAnonymous(ctxReq context.Context) <-chan AccessTokenResponse {
	result := make(chan AccessTokenResponse)
	go func() {
		defer close(result)
		claims := Claim{}
		claims.Issuer = model.Bhinneka
		claims.DeviceID = model.DefaultDeviceID
		claims.DeviceLogin = model.DefaultDeviceLogin
		claims.Audience = "sturgeon-generator"
		claims.Subject = model.DefaultSubject
		claims.Authorised = false

		// redis key format: STG-USR123-ASX1234-WEB
		redisLoginKey := strings.Join([]string{"STG", claims.Subject, claims.DeviceID, claims.DeviceLogin}, "-")

		loadResult := <-j.LoginSessionRepo.Load(ctxReq, redisLoginKey)
		if loadResult.Error == nil {
			existingToken, _ := loadResult.Result.(model.LoginSessionRedis)
			result <- AccessTokenResponse{Error: nil, AccessToken: AccessToken{AccessToken: "Bearer " + existingToken.Token}}
			return
		}

		// generate token based on claims
		tokenResult := <-j.GenerateAccessToken(claims)
		token := tokenResult.AccessToken.AccessToken
		expiredAt := tokenResult.AccessToken.ExpiredAt
		jti := tokenResult.AccessToken.JTI

		paramRedis := &model.LoginSessionRedis{
			Key:         redisLoginKey,
			Token:       token,
			ExpiredTime: expiredAt.Sub(time.Now()),
		}

		j.LoginSessionRepo.Save(ctxReq, paramRedis)

		result <- AccessTokenResponse{Error: nil, AccessToken: AccessToken{AccessToken: "Bearer " + token, ExpiredAt: expiredAt, JTI: jti}}
	}()

	return result

}
