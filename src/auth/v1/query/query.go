package query

import "context"

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// AuthQuery interface abstraction
type AuthQuery interface {
	UpdateLastLogin(uid string) <-chan ResultQuery
	GetAccountId(ctxReq context.Context, contactId int) <-chan ResultQuery
}

// AuthQueryOA interface abstraction
type AuthQueryOA interface {
	GetAzureToken(ctxReq context.Context, code, redirectURI string) <-chan ResultQuery
	GetGoogleToken(ctxReq context.Context, code, redirectURI string) <-chan ResultQuery
	GetGoogleTokenInfo(ctxReq context.Context, token string) <-chan ResultQuery
	GetFacebookToken(ctxReq context.Context, code, redirectURI string) <-chan ResultQuery
	GetAppleToken(ctxReq context.Context, code, redirectURI, clientID string) <-chan ResultQuery
	GetDetailAzureMember(ctxReq context.Context, token string) <-chan ResultQuery
	GetDetailFacebookMember(ctxReq context.Context, code string) <-chan ResultQuery
	GetDetailGoogleMember(ctxReq context.Context, code string) <-chan ResultQuery
}
