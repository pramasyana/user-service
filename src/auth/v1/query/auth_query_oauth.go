package query

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/labstack/echo"

	"strings"

	"github.com/Bhinneka/golib/tracer"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/config/rsa"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

const (
	grantType         = "grant_type"
	textclientID      = "client_id"
	clientSecret      = "client_secret"
	redirectURIKey    = "redirect_uri"
	authorizationCode = "authorization_code"
	tagCode           = "code"
)

// AuthQueryOAuth data structure
type AuthQueryOAuth struct {
	FacebookBaseURL      *url.URL
	FacebookTokenBaseURL *url.URL
	GoogleBaseURL        *url.URL
	GoogleTokenBaseURL   *url.URL
	GoogleOAuthBaseURL   *url.URL
	AzureBaseURL         *url.URL
	AzureLoginBaseURL    *url.URL
	AppleBaseURL         *url.URL
	GoogleClientID       string
	GoogleClientSecret   string
	FacebookClientID     string
	FacebookClientSecret string
	AzureADTenanID       string
	AzureADClientID      string
	AzureADClientSecret  string
	AzureADResource      string
	AppleClientID        string
	AppleTeamID          string
	AppleKeyID           string
}

// NewAuthQueryOAuth function for initializing member query
func NewAuthQueryOAuth(cfg localConfig.OAuthService) *AuthQueryOAuth {
	return &AuthQueryOAuth{
		FacebookBaseURL:      cfg.FacebookBaseURL,
		FacebookTokenBaseURL: cfg.FacebookTokenBaseURL,
		GoogleBaseURL:        cfg.GoogleBaseURL,
		GoogleTokenBaseURL:   cfg.GoogleBaseTokenURL,
		GoogleOAuthBaseURL:   cfg.GoogleOAuthBaseURL,
		AzureBaseURL:         cfg.AzureBaseURL,
		AzureLoginBaseURL:    cfg.AzureLoginBaseURL,
		AppleBaseURL:         cfg.AppleBaseURL,
		FacebookClientID:     cfg.FBClientID,
		FacebookClientSecret: cfg.FBClientSecret,
		GoogleClientID:       cfg.GoogleClientID,
		GoogleClientSecret:   cfg.GoogleClientSecret,
		AzureADTenanID:       cfg.AzureADTenanID,
		AzureADClientID:      cfg.AzureADClientID,
		AzureADClientSecret:  cfg.AzureADClientSecret,
		AzureADResource:      cfg.AzureADResource,
		AppleTeamID:          cfg.AppleTeamID,
		AppleKeyID:           cfg.AppleKeyID,
	}
}

// GetAzureToken function for getting azure token
func (oa *AuthQueryOAuth) GetAzureToken(ctxReq context.Context, code, redirectURI string) <-chan ResultQuery {
	ctx := "AuthQuery-GetAzureToken"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[tagCode] = code

		uri := fmt.Sprintf("%s/%s/oauth2/token", oa.AzureLoginBaseURL.String(), oa.AzureADTenanID)

		form := url.Values{}
		form.Add(grantType, authorizationCode)
		form.Add(tagCode, code)
		form.Add(textclientID, oa.AzureADClientID)
		form.Add(clientSecret, oa.AzureADClientSecret)
		form.Add("resource", oa.AzureADResource)
		form.Add(redirectURIKey, redirectURI)

		headers := map[string]string{
			echo.HeaderContentType: echo.MIMEApplicationForm,
		}

		azureToken := model.AuthAzureToken{}

		if err := helper.GetHTTPNewRequest(ctxReq, "POST", uri, strings.NewReader(form.Encode()), &azureToken, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_azure_token", err, form)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		// prevent error
		if len(azureToken.Error) > 0 {
			err := errors.New(azureToken.ErrorDescription)
			helper.SendErrorLog(ctxReq, ctx, "error_azure_token", err, form)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: azureToken}
	})

	return output
}

// GetGoogleToken function for getting google token
func (oa *AuthQueryOAuth) GetGoogleToken(ctxReq context.Context, code, redirectURI string) <-chan ResultQuery {
	ctx := "AuthQuery-GetGoogleToken"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[tagCode] = code

		uri := fmt.Sprintf("%s/o/oauth2/token", oa.GoogleTokenBaseURL.String())

		form := url.Values{}
		form.Add(grantType, authorizationCode)
		form.Add(tagCode, code)
		form.Add(textclientID, oa.GoogleClientID)
		form.Add(clientSecret, oa.GoogleClientSecret)
		form.Add(redirectURIKey, redirectURI)

		headers := map[string]string{
			echo.HeaderContentType: echo.MIMEApplicationForm,
		}

		googleToken := model.AuthGoogleToken{}

		if err := helper.GetHTTPNewRequest(ctxReq, "POST", uri, strings.NewReader(form.Encode()), &googleToken, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_google_token", err, form)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		// prevent error
		if len(googleToken.Error) > 0 {
			helper.SendErrorLog(ctxReq, ctx, "GetGoogleTokenCheckLen", errors.New(googleToken.Error), form)
			err := errors.New(googleToken.ErrorDescription)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: googleToken}
	})

	return output
}

// GetGoogleTokenInfo function for getting google token Info
func (oa *AuthQueryOAuth) GetGoogleTokenInfo(ctxReq context.Context, token string) <-chan ResultQuery {
	ctx := "AuthQuery-GetGoogleTokenInfo"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags["token"] = token

		uri := fmt.Sprintf("%s/tokeninfo?id_token=%s", oa.GoogleOAuthBaseURL.String(), token)
		googleToken := model.GoogleOAuthToken{}

		if err := helper.GetHTTPNewRequest(ctxReq, "POST", uri, nil, &googleToken, nil); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_google_token", err, nil)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: googleToken}
	})

	return output
}

// GetFacebookToken function for getting facebook token
func (oa *AuthQueryOAuth) GetFacebookToken(ctxReq context.Context, code, redirectURI string) <-chan ResultQuery {
	ctx := "AuthQuery-GetFacebookToken"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[tagCode] = code

		uri := fmt.Sprintf("%s/oauth/access_token", oa.FacebookTokenBaseURL.String())

		form := url.Values{}
		form.Add(grantType, authorizationCode)
		form.Add(tagCode, code)
		form.Add(textclientID, oa.FacebookClientID)
		form.Add(clientSecret, oa.FacebookClientSecret)
		form.Add(redirectURIKey, redirectURI)

		headers := map[string]string{
			echo.HeaderContentType: echo.MIMEApplicationForm,
		}

		facebookToken := model.AuthFacebookToken{}

		if err := helper.GetHTTPNewRequest(ctxReq, "POST", uri, strings.NewReader(form.Encode()), &facebookToken, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_facebook_token", err, form)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		// prevent error
		if facebookToken.Error != nil {
			err := errors.New(facebookToken.Error.Message)
			helper.SendErrorLog(ctxReq, ctx, "error_facebook_token", err, form)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: facebookToken}
	})

	return output
}

// GetDetailAzureMember function for getting detail azure member by code
func (oa *AuthQueryOAuth) GetDetailAzureMember(ctxReq context.Context, token string) <-chan ResultQuery {
	ctx := "MemberQuery-GetDetailAzureMember"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[tagCode] = token

		headers := map[string]string{
			echo.HeaderContentType:   echo.MIMEApplicationJSON,
			echo.HeaderAuthorization: fmt.Sprintf("Bearer %s", token),
		}
		uri := fmt.Sprintf("%s/%s/me?api-version=1.6", oa.AzureBaseURL.String(), oa.AzureADTenanID)

		azure := model.AzureResponse{}

		if err := helper.GetHTTPNewRequest(ctxReq, "GET", uri, nil, &azure, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_azure_data", err, nil)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: azure}
	})

	return output
}

// GetDetailFacebookMember function for getting detail facebook member by code
func (oa *AuthQueryOAuth) GetDetailFacebookMember(ctxReq context.Context, code string) <-chan ResultQuery {
	ctx := "MemberQuery-GetDetailFacebookMember"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[tagCode] = code
		uri := fmt.Sprintf("%s/me?fields=id,name,first_name,last_name,email,birthday,gender&access_token=%s", oa.FacebookBaseURL.String(), code)

		facebook := model.FacebookResponse{}

		if err := helper.GetHTTPNewRequest(ctxReq, "GET", uri, nil, &facebook); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_facebook_data", err, nil)
			errC := errors.New("failed to get data from facebook")
			output <- ResultQuery{Error: errC}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: facebook}
	})

	return output
}

// GetDetailGoogleMember function for getting detail google member by code
func (oa *AuthQueryOAuth) GetDetailGoogleMember(ctxReq context.Context, code string) <-chan ResultQuery {
	ctx := "MemberQuery-GetDetailGoogleMember"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[tagCode] = code
		uri := fmt.Sprintf("%s/oauth2/v2/userinfo", oa.GoogleBaseURL.String())

		headers := map[string]string{
			echo.HeaderAuthorization: fmt.Sprintf("Bearer %s", code),
		}
		google := model.GoogleOAuth2Response{}

		if err := helper.GetHTTPNewRequest(ctxReq, "GET", uri, nil, &google, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_google_data", err, nil)
			errC := errors.New("failed to get data from google")
			output <- ResultQuery{Error: errC}
			tags[helper.TextResponse] = err
			return
		}

		if google.Error != nil {
			err := errors.New("failed to get data from google, check the token you provide")
			helper.SendErrorLog(ctxReq, ctx, "get_google_data_error", err, google.Error)
			output <- ResultQuery{Error: err}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: google}
	})

	return output
}

// GetAppleToken function for getting detail apple member by code & clientSecretCode
func (oa *AuthQueryOAuth) GetAppleToken(ctxReq context.Context, code, redirectURI, clientID string) <-chan ResultQuery {
	ctx := "MemberQuery-GetAppleToken"
	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		uri := fmt.Sprintf("%s/auth/token", oa.AppleBaseURL.String())
		tags[helper.TextParameter] = uri
		tags[tagCode] = code

		clientSecretCode, err := rsa.InitAppleClientSecret(clientID, oa.AppleTeamID, oa.AppleKeyID)
		if err != nil {
			output <- ResultQuery{Error: errors.New("failed to get data from apple")}
			helper.SendErrorLog(ctxReq, ctx, "get_apple_token", err, nil)
			return
		}

		form := url.Values{}
		form.Add(grantType, authorizationCode)
		form.Add(tagCode, code)
		form.Add(textclientID, clientID)
		form.Add(clientSecret, clientSecretCode)
		form.Add(redirectURIKey, redirectURI)

		headers := map[string]string{
			echo.HeaderContentType: echo.MIMEApplicationForm,
		}

		appleResponse := model.AppleResponse{}

		if err := helper.GetHTTPNewRequestV2(ctxReq, "POST", uri, strings.NewReader(form.Encode()), &appleResponse, headers); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_apple_data", err, form)
			output <- ResultQuery{Error: errors.New("failed to get data from apple")}
			tags[helper.TextResponse] = err
			return
		}

		output <- ResultQuery{Result: appleResponse}
	})
	return output
}
