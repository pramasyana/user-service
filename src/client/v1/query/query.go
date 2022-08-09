package query

import "context"

type ClientQuery interface {
	Validate(ctxReq context.Context, clientID, clientSecret string) (bool, error)
}
