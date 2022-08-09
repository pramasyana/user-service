package config

import (
	"net/url"
	"time"

	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	authServices "github.com/Bhinneka/user-service/src/auth/v1/service"
	authToken "github.com/Bhinneka/user-service/src/auth/v1/token"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	corporateRepo "github.com/Bhinneka/user-service/src/corporate/v2/repo"
	documentRepo "github.com/Bhinneka/user-service/src/document/v2/repo"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	paymentsRepo "github.com/Bhinneka/user-service/src/payments/v1/repo"
	"github.com/Bhinneka/user-service/src/service"
	sessionInfoQuery "github.com/Bhinneka/user-service/src/session/v1/query"
	sessionInfoRepo "github.com/Bhinneka/user-service/src/session/v1/repo"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
	shippingAddressRepo "github.com/Bhinneka/user-service/src/shipping_address/v2/repo"
)

// ServiceRepository general parameter
type ServiceRepository struct {
	ClientAppRepoRead                  authRepo.ClientAppRepository
	ClientAppRepoWrite                 authRepo.ClientAppRepository
	CorporateAccountRepository         corporateRepo.AccountRepository
	CorporateAccountTempRepository     corporateRepo.AccountTemporaryRepository
	CorporateAccountContactRepository  corporateRepo.AccountContactRepository
	CorporateContactRepository         corporateRepo.ContactRepository
	CorporateAddressRepository         corporateRepo.AddressRepository
	CorporatePhoneRepository           corporateRepo.PhoneRepository
	CorporateDocumentRepository        corporateRepo.DocumentRepository
	CorporateContactNPWPRepository     corporateRepo.ContactNpwpRepository
	CorporateContactAddressRepository  corporateRepo.ContactAddressRepository
	CorporateContactTempRepository     corporateRepo.ContactTempRepository
	CorporateLeadsRepository           corporateRepo.LeadsRepository
	MerchantRepository                 merchantRepo.MerchantRepository
	MerchantDocumentRepository         merchantRepo.MerchantDocumentRepository
	MerchantBankRepository             merchantRepo.MerchantBankRepository
	MerchantEmployeeRepository         merchantRepo.MerchantEmployeeRepository
	MerchantAddressRepository          merchantRepo.MerchantAddressRepository
	ShippingAddressRepository          shippingAddressRepo.ShippingAddressRepository
	ShippingAddressRedisRepository     shippingAddressRepo.ShippingAddressRepositoryRedis
	MemberRepository                   memberRepo.MemberRepository
	MemberMFARepository                memberRepo.MemberMFARepository
	MemberRedisRepository              memberRepo.MemberRepositoryRedis
	TokenActivationRepoRedis           memberRepo.TokenActivationRepository
	AttemptRepositoryRedis             authRepo.AttemptRepository
	LoginSessionRepositoryRedis        authRepo.LoginSessionRepository
	RefreshTokenRepository             authRepo.RefreshTokenRepository
	SessionInfoRepo                    sessionInfoRepo.SessionInfoRepository
	DocumentRepository                 documentRepo.DocumentRepository
	Repository                         *sharedRepo.Repository
	CorporateContactDocumentRepository corporateRepo.ContactDocumentRepository
	PaymentsRepository                 paymentsRepo.PaymentsRepository
}

// ServiceQuery general parameter
type ServiceQuery struct {
	CorporateContactQueryRead    corporateQuery.ContactQuery
	CorporateAccContactQueryRead corporateQuery.AccountContactQuery
	MemberQueryRead              memberQuery.MemberQuery
	MemberQueryWrite             memberQuery.MemberQuery
	MemberMFAQueryRead           memberQuery.MemberMFAQuery
	SessionInfoQuery             sessionInfoQuery.SessionInfoQuery
}

// ServiceShared general parameter
type ServiceShared struct {
	StaticService       service.StaticServices
	UploadService       service.UploadServices
	MerchantService     service.MerchantServices
	ActivityService     service.ActivityServices
	BarracudaService    service.BarracudaServices
	QPublisher          service.QPublisher
	NotificationService service.NotificationServices
	SendbirdService     service.SendbirdServices
}

// OAuthService general struct
type OAuthService struct {
	AzureLoginBaseURL    *url.URL
	AzureBaseURL         *url.URL
	FacebookBaseURL      *url.URL
	FacebookTokenBaseURL *url.URL
	GoogleBaseURL        *url.URL
	GoogleBaseTokenURL   *url.URL
	GoogleOAuthBaseURL   *url.URL
	AppleBaseURL         *url.URL
	FBClientID           string
	FBClientSecret       string
	GoogleClientID       string
	GoogleClientSecret   string
	AzureADTenanID       string
	AzureADClientID      string
	AzureADClientSecret  string
	AzureADResource      string
	AppleTeamID          string
	AppleKeyID           string
}

// MembershipParameters member parameter
type MembershipParameters struct {
	Hash                              memberModel.PasswordHasher
	TokenActivationExpiration         time.Duration
	ResendActivationAttemptAge        string
	ResendActivationAttemptAgeRequest string
	Topic                             string
	IsProductionStage                 bool
	SturgeonCFUrl                     string
	B2cCFUrl                          string
	AccessTokenGenerator              authToken.AccessTokenGenerator
}

// AuthParameters auth parameter
type AuthParameters struct {
	Hash                   memberModel.PasswordHasher
	AuthServices           authServices.LDAPService
	GoogleVerifyCaptchaURL *url.URL
	AccessTokenGenerator   authToken.AccessTokenGenerator
	RefreshTokenAge        string
	B2cCFUrl               string
	LoginAttemptAge        string
	Topic                  string
	IsProductionStage      bool
	SpecialRefreshTokenAge string
	EmailSpecialTokenAge   string
}
