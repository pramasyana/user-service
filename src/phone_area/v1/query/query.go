package query

import "context"

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// PhoneAreaQuery interface abstraction
type PhoneAreaQuery interface {
	FindAll(ctxReq context.Context) <-chan ResultQuery
	Count(ctxReq context.Context) <-chan ResultQuery
}
