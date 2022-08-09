package usecase

import (
	"context"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

var findEmail = []string{"##YEAR##", "##FULLNAME##", "##URL##"}

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	Error      error
	HTTPStatus int
	ErrorData  []model.MemberError
}

// MemberUseCase interface abstraction
type MemberUseCase interface {
	CheckEmailAndMobileExistence(ctxReq context.Context, data *model.Member) <-chan ResultUseCase
	GetDetailMemberByEmail(email string) <-chan ResultUseCase
	GetDetailMemberByID(ctxReq context.Context, uid string) <-chan ResultUseCase
	UpdateDetailMemberByID(ctxReq context.Context, data model.Member) <-chan ResultUseCase
	UpdatePassword(ctxReq context.Context, token, uid, oldPassword, newPassword string) <-chan ResultUseCase
	AddNewPassword(ctxReq context.Context, data model.Member) <-chan ResultUseCase
	RegisterMember(ctxReq context.Context, data *model.Member) <-chan ResultUseCase
	ActivateMember(ctxReq context.Context, token, requestFrom string) <-chan ResultUseCase
	ForgotPassword(ctxReq context.Context, email string) <-chan ResultUseCase
	ValidateToken(ctxReq context.Context, token string) <-chan ResultUseCase

	ChangeForgotPassword(ctxReq context.Context, token, newPassword, rePassword, requestFrom string) <-chan ResultUseCase

	ActivateNewPassword(ctxReq context.Context, token, password, rePassword string) <-chan ResultUseCase
	GetListMembers(ctxReq context.Context, params *model.Parameters) <-chan ResultUseCase
	RegenerateToken(ctxReq context.Context, data model.Member) <-chan ResultUseCase
	MigrateMember(ctxReq context.Context, members *model.Members) <-chan ResultUseCase
	ResendActivation(ctxReq context.Context, email string) <-chan ResultUseCase

	PublishToKafkaUser(ctxReq context.Context, data *model.Member, eventType string) error
	ValidateEmailDomain(ctxReq context.Context, email string) <-chan ResultUseCase
	UpdateProfilePicture(ctxReq context.Context, data model.ProfilePicture) <-chan ResultUseCase

	RevokeAllAccess(ctxReq context.Context, uid string, token string) <-chan ResultUseCase
	GetLoginActivity(ctxReq context.Context, params *model.ParametersLoginActivity) <-chan ResultUseCase
	GetProfileComplete(ctxReq context.Context, uid string) <-chan ResultUseCase
	UpdateProfileName(ctxReq context.Context, data model.ProfileName) <-chan ResultUseCase
	RevokeAccess(ctxReq context.Context, uid string, jti string) <-chan ResultUseCase

	// MFA Related
	GetMFASettings(ctxReq context.Context, uid string) <-chan ResultUseCase
	ActivateMFASettings(ctxReq context.Context, activateData model.MFAActivateSettings) <-chan ResultUseCase
	ActivateMFASettingV3(ctxReq context.Context, activateData model.MFAActivateSettings) <-chan ResultUseCase

	// Import Related
	ParseMemberData(ctxReq context.Context, input []byte) ([]*model.Member, error)
	ValidateEmailAndPhone(ctxReq context.Context, data *model.Member) error
	BulkValidateEmailAndPhone(ctxReq context.Context, members []*model.Member) []string

	// Narwhal MFA
	GetNarwhalMFASettings(ctxReq context.Context, uid string) <-chan ResultUseCase
	// both side
	GenerateMFASettings(ctxReq context.Context, userID, requestFrom string) <-chan ResultUseCase
	DisabledMFASetting(ctxReq context.Context, userID, requestFrom string) <-chan ResultUseCase
	ImportMember(ctxReq context.Context, data *model.Member) <-chan ResultUseCase
	BulkImportMember(ctxReq context.Context, data []*model.Member) <-chan ResultUseCase

	// Email related
	SendEmailRegisterMember(ctxReq context.Context, data model.SuccessResponse, registerType string) <-chan ResultUseCase
	SendEmailWelcomeMember(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase
	SendEmailForgotPassword(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase
	SendEmailSuccessForgotPassword(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase
	SendEmailAddMember(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase

	// Log related
	InsertLogMember(ctxReq context.Context, oldData, newData *model.Member, action string) error

	// get session token from sendbird
	GetSendbirdToken(ctxReq context.Context, params *serviceModel.SendbirdRequest) ResultUseCase
	GetSendbirdTokenV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) ResultUseCase
	CheckSendbirdToken(ctxReq context.Context, params *serviceModel.SendbirdRequest) ResultUseCase
	CheckSendbirdTokenV4(ctxReq context.Context, params *serviceModel.SendbirdRequestV4) ResultUseCase

	// sync password
	SyncPassword(ctxReq context.Context, token, oldPassword, newPassword string) <-chan ResultUseCase
	ChangePassword(ctxReq context.Context, token, oldPassword, newPassword string) <-chan ResultUseCase
	Clients(ctxReq context.Context, token string) <-chan ResultUseCase

	// employee
	ActivateMerchantEmployee(ctxReq context.Context, token string) <-chan ResultUseCase
}
