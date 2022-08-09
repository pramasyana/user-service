package service

import (
	"context"

	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

//QPublisher interface, publisher interface abstraction
type QPublisher interface {
	Publish(ctxReq context.Context, topic string, messageKey shared.MessageKey, message []byte) error
	PublishKafka(ctxReq context.Context, topic string, messageKey string, message []byte) error
	QueueJob(ctxReq context.Context, payload interface{}, messageKey, jobType string) error
	BulkPublishKafka(ctxReq context.Context, topic string, messages []serviceModel.Messages) error
}

//StaticServices interface, publisher interface abstraction
type StaticServices interface {
	FindStaticsByID(ctxReq context.Context, id string) <-chan serviceModel.ServiceResult
	FindStaticsGwsByID(ctxReq context.Context, id string) <-chan serviceModel.ServiceResult
}

type NotificationServices interface {
	GetTemplateByID(ctxReq context.Context, templateId, envKey string) <-chan serviceModel.ServiceResult
	SendEmail(ctxReq context.Context, email serviceModel.Email) (string, error)
}

//MerchantServices interface, publisher interface abstraction
type MerchantServices interface {
	FindMerchantServiceByID(ctxReq context.Context, id, token, merchantID string) <-chan serviceModel.ServiceResult
	PublishToKafkaUserMerchant(ctxReq context.Context, data *model.B2CMerchantDataV2, eventType, producer string) error
	InsertLogMerchant(ctxReq context.Context, oldData, newData model.B2CMerchantDataV2, action, module string) error
	InsertLogMerchantPIC(ctxReq context.Context, oldData, newData model.B2CMerchantDataV2, action, module string, member memberModel.Member) error
}

//BarracudaServices interface, publisher interface abstraction
type BarracudaServices interface {
	FindZipcode(ctxReq context.Context, data serviceModel.ZipCodeQueryParameter) <-chan serviceModel.ServiceResult
}

//UploadServices interface, publisher interface abstraction
type UploadServices interface {
	GetURLImage(ctxReq context.Context, url string, isAttachment string) <-chan serviceModel.ServiceResult
}

//ActivityServices interface, publisher interface abstraction
type ActivityServices interface {
	InsertLog(ctxReq context.Context, oldData, newData interface{}, payload serviceModel.Payload) error
	CreateLog(ctxReq context.Context, param serviceModel.Payload) <-chan serviceModel.ServiceResult
	GetAll(ctxReq context.Context, param *sharedModel.Parameters) <-chan serviceModel.ServiceResult
	GetLogByID(ctxReq context.Context, logID string) <-chan serviceModel.ServiceResult
}

// SendbirdServices interface, publisher interface abstraction
type SendbirdServices interface {
	CheckUserSenbird(ctxReq context.Context, data *serviceModel.SendbirdRequest) serviceModel.ServiceResult
	CheckUserSenbirdV4(ctxReq context.Context, data *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult
	UpdateUserSendbird(ctxReq context.Context, data *serviceModel.SendbirdRequest) serviceModel.ServiceResult
	UpdateUserSendbirdV4(ctxReq context.Context, data *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult
	CreateUserSendbird(ctxReq context.Context, data *serviceModel.SendbirdRequest) serviceModel.ServiceResult
	CreateUserSendbirdV4(ctxReq context.Context, data *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult
	GetTokenUserSendbird(ctxReq context.Context, data *serviceModel.SendbirdRequest) serviceModel.ServiceResult
	CreateTokenUserSendbird(ctxReq context.Context, data *serviceModel.SendbirdRequest) serviceModel.ServiceResult
	CreateTokenUserSendbirdV4(ctxReq context.Context, data *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult
	GetUserSendbird(ctxReq context.Context, data *serviceModel.SendbirdRequest) serviceModel.ServiceResult
	GetUserSendbirdV4(ctxReq context.Context, data *serviceModel.SendbirdRequestV4) serviceModel.ServiceResult
}
