package model

import "strings"

const (
	PARAM_DOMAIN_BHINNEKA = "bhinneka.com"
)

// FacebookResponse data structure
type FacebookResponse struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Birthday  string      `json:"birthday"`
	Gender    string      `json:"gender"`
	Error     interface{} `json:"error"`
	LastName  string      `json:"last_name"`
	FirstName string      `json:"first_name"`
}

// GoogleOAuth2Response google data structure
type GoogleOAuth2Response struct {
	FamilyName string      `json:"family_name"`
	GivenName  string      `json:"given_name"`
	Name       string      `json:"name"`
	Picture    string      `json:"picture"`
	Gender     string      `json:"gender"`
	Email      string      `json:"email"`
	ID         string      `json:"id"`
	HD         string      `json:"hd"`
	Error      *GooglError `json:"error"`
}

// GooglError data structure
type GooglError struct {
	Errors  interface{} `json:"errors"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
}

// AuthAzureToken data structure to collect token from azure
type AuthAzureToken struct {
	Token            string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// AuthGoogleToken data structure to collect token from google
type AuthGoogleToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	IDToken          string `json:"id_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// AuthFacebookToken data structure to collect token from facebook
type AuthFacebookToken struct {
	AccessToken string         `json:"access_token"`
	TokenType   string         `json:"token_type"`
	ExpiresIn   int            `json:"expires_in"`
	Error       *FacebookError `json:"error,omitempty"`
}

// FacebookError error
type FacebookError struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	Code      int    `json:"code"`
	FbtraceID string `json:"fbtrace_id"`
}

// AzureResponse data structure
type AzureResponse struct {
	ObjectID         string `json:"objectId"`
	CompanyName      string `json:"companyName"`
	Email            string `json:"mail"`
	MailNickname     string `json:"mailNickname"`
	DisplayName      string `json:"displayName"`
	Department       string `json:"department"`
	JobTitle         string `json:"jobTitle"`
	ImmutableID      string `json:"immutableID"`
	Error            string `json:"error"`
	ErrorDescription string `json:"errorDescription"`
}

// LDAPProfile response login data struct
type LDAPProfile struct {
	Email       string `json:"mail"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
	Department  string `json:"department"`
	ObjectID    string `json:"objectId"`
	JobTitle    string `json:"jobTitle"`
}

type GoogleOAuthToken struct {
	Iss           string `json:"iss"`
	Nbf           string `json:"nbf"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Azp           string `json:"azp"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Iat           string `json:"iat"`
	Exp           string `json:"exp"`
	Jti           string `json:"jti"`
}

// IsBhinnekaEmail function for check valid bhinneka's email
func (a AzureResponse) IsBhinnekaEmail() bool {
	bhinnekaDomain := "@" + PARAM_DOMAIN_BHINNEKA
	return strings.Contains(a.Email, bhinnekaDomain)
}

// IsBhinnekaEmail function for check valid bhinneka's email
func (g GoogleOAuth2Response) IsBhinnekaEmail() bool {
	bhinnekaDomain := PARAM_DOMAIN_BHINNEKA
	return len(g.HD) > 0 && g.HD == bhinnekaDomain
}

// IsBhinnekaEmail function for check valid bhinneka's email
func (g GoogleOAuthToken) IsBhinnekaEmail() bool {
	bhinnekaDomain := PARAM_DOMAIN_BHINNEKA
	return len(g.Hd) > 0 && g.Hd == bhinnekaDomain
}

// IsBhinnekaEmail function for check valid bhinneka's email
func (a AppleProfile) IsBhinnekaEmail() bool {
	bhinnekaDomain := PARAM_DOMAIN_BHINNEKA
	return strings.Contains(a.Email, bhinnekaDomain)
}

// AppleResponse apple data structure
type AppleResponse struct {
	IDToken string `json:"id_token"`
	Error   string `json:"error"`
}

// AppleProfile apple data structure
type AppleProfile struct {
	Sub            string `json:"sub"`
	Email          string `json:"id_token"`
	IsPrivateEmail string `json:"is_private_email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
}

// MicrositeClient base struct
type MicrositeClient struct {
	Firstname string
	Email     string
}
