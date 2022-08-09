package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// ClientAppRepository auth data abstraction
type ClientAppRepository interface {
	Save(*model.ClientApp) <-chan ResultRepository
	Load(int) <-chan ResultRepository
	FindByClientID(string) <-chan ResultRepository
}

// RefreshTokenRepository interface abstraction
type RefreshTokenRepository interface {
	Save(ctxReq context.Context, refreshToken *model.RefreshToken) <-chan ResultRepository
	Load(ctxReq context.Context, key string) <-chan ResultRepository
	Delete(ctxReq context.Context, key string) <-chan ResultRepository
}

// AttemptRepository interface abstraction
type AttemptRepository interface {
	Save(ctxReq context.Context, data *model.LoginAttempt) <-chan ResultRepository
	Load(ctxReq context.Context, key string) <-chan ResultRepository
	Delete(ctxReq context.Context, key string) <-chan ResultRepository
}

// LoginSessionRepository interface abstraction
type LoginSessionRepository interface {
	Save(ctxReq context.Context, data *model.LoginSessionRedis) <-chan ResultRepository
	Load(ctxReq context.Context, key string) <-chan ResultRepository
	Delete(ctxReq context.Context, key string) <-chan ResultRepository
	GetLoginActive(ctxReq context.Context, key string) <-chan ResultRepository
}
