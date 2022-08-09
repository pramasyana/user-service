package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// MemberRepository interface abstraction
type MemberRepository interface {
	Save(ctxReq context.Context, member model.Member) <-chan ResultRepository
	Load(ctxReq context.Context, uid string) <-chan ResultRepository
	LoadMember(uid string) ResultRepository
	BulkSave(ctxReq context.Context, member []model.Member) <-chan ResultRepository
	BulkImportSave(ctxReq context.Context, member []*model.Member) <-chan ResultRepository
	FindMaxID(ctxReq context.Context) <-chan ResultRepository
	UpdateProfilePicture(ctxReq context.Context, data model.ProfilePicture) <-chan ResultRepository
	UpdateFlagIsSyncMember(ctxReq context.Context, member model.Member) <-chan ResultRepository
	UpdatePasswordMemberByEmail(ctxReq context.Context, member model.Member) <-chan ResultRepository
}

// MemberRepositoryRedis interface abstraction
type MemberRepositoryRedis interface {
	Save(stc *model.MemberRedis) <-chan ResultRepository
	Load(uid string) <-chan ResultRepository
	Delete(uid string) <-chan ResultRepository
	LoadByKey(ctxReq context.Context, key string) <-chan ResultRepository
	SaveResendActivationAttempt(ctxReq context.Context, data *model.ResendActivationAttempt) <-chan ResultRepository
	RevokeAllAccess(ctxReq context.Context, key, activeKey string) <-chan ResultRepository
}

//TokenActivationRepository interface abstraction
type TokenActivationRepository interface {
	Save(*model.TokenActivation) <-chan ResultRepository
	Load(key string) <-chan ResultRepository
	Delete(key string) <-chan ResultRepository
}

// DolphinLogRepository interface
type DolphinLogRepository interface {
	Save(context.Context, *model.DolphinLog) error
	Load(int) ResultRepository
}

// MemberAdditionalInfoRepository interface
type MemberAdditionalInfoRepository interface {
	Save(ctxReq context.Context, data *model.MemberAdditionalInfo) <-chan ResultRepository
	Update(ctxReq context.Context, data *model.MemberAdditionalInfo) <-chan ResultRepository
	Load(ctxReq context.Context, uid string, authType string) <-chan ResultRepository
}

// MemberMFARepository interface
type MemberMFARepository interface {
	MFAEnabled(ctxReq context.Context, uid string, mfaKey string) <-chan ResultRepository
	MFADisabled(ctxReq context.Context, uid string) <-chan ResultRepository
	EnableNarwhalMFA(ctxReq context.Context, uid string, mfaKey string) <-chan ResultRepository
	DisableNarwhalMFA(ctxReq context.Context, uid string) <-chan ResultRepository
}
