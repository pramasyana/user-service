package usecase

import (
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	"github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
)

const (
	errMsgFailedSendEmail = "failed to send email"
	msgErrorSave          = "failed to save"
	msgErrorFindAddress   = "failed to find address"
	msgErrorUpdatePrimary = "failed to update primary address"
	msgErrorDeleteAddress = "failed to delete address"
)

// MerchantUseCaseImpl data structure
type MerchantUseCaseImpl struct {
	Repository           *sharedRepo.Repository
	MerchantRepo         repo.MerchantRepository
	MerchantAddressRepo  repo.MerchantAddressRepository
	MerchantBankRepo     repo.MerchantBankRepository
	MerchantEmployeeRepo repo.MerchantEmployeeRepository
	MerchantDocumentRepo repo.MerchantDocumentRepository
	MemberRepoRead       memberRepo.MemberRepository
	UploadService        service.UploadServices
	MerchantService      service.MerchantServices
	TokenGenerator       token.AccessTokenGenerator
	MemberQueryRead      query.MemberQuery
	NotificationService  service.NotificationServices
	QueuePublisher       service.QPublisher
	SendbirdService      service.SendbirdServices
}

// NewMerchantUseCase function for initialise merchant use case implementation mo el
func NewMerchantUseCase(repository localConfig.ServiceRepository, services localConfig.ServiceShared, tokenGenerator token.AccessTokenGenerator, query localConfig.ServiceQuery) MerchantUseCase {
	return &MerchantUseCaseImpl{
		Repository:           repository.Repository,
		MerchantRepo:         repository.MerchantRepository,
		MerchantAddressRepo:  repository.MerchantAddressRepository,
		MerchantBankRepo:     repository.MerchantBankRepository,
		MerchantEmployeeRepo: repository.MerchantEmployeeRepository,
		MerchantDocumentRepo: repository.MerchantDocumentRepository,
		MemberRepoRead:       repository.MemberRepository,
		UploadService:        services.UploadService,
		MerchantService:      services.MerchantService,
		TokenGenerator:       tokenGenerator,
		MemberQueryRead:      query.MemberQueryRead,
		NotificationService:  services.NotificationService,
		QueuePublisher:       services.QPublisher,
		SendbirdService:      services.SendbirdService,
	}
}
