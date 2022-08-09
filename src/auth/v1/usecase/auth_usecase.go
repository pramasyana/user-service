package usecase

import (
	"net/url"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/src/auth/v1/query"
	"github.com/Bhinneka/user-service/src/auth/v1/repo"
	authServices "github.com/Bhinneka/user-service/src/auth/v1/service"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantRepoRead "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	sessionInfoRepo "github.com/Bhinneka/user-service/src/session/v1/repo"
)

const (
	scopeErrorParseSocmed             = "error_parse_socmed_data"
	scopeParseLdap                    = "parse_ldap"
	scopeValidateRequestTokenFacebook = "validate_request_token_facebook"
	scopeGetPublicKey                 = "get_public_key"
	scopeValidateToken                = "validate_token"
	msgErrorEmailRegisterAccount      = "cannot register your account, your email is hidden"
	msgResultNotMember                = "result is not member"
	msgResultNotAccount               = "result is not account"
	msgFailedLoginSocmedEmail         = "failed to login, email doesn't match with your session email"
	msgFailedToSendEmail              = "failed to send email"
	keyAttempt                        = "ATTEMPT:%s"
	msgTokenExpired                   = "this token has been expired"
	textGetTokenPass                  = "auth_get_token_pass_params"
	textGetToken                      = "auth_get_token_params"
	signUpFromLKPP                    = "lkpp"
	msgUserIdNotValid                 = "user id is not valid"
)

// AuthUseCaseImpl data structure
type AuthUseCaseImpl struct {
	ClientAppRepoRead            repo.ClientAppRepository
	ClientAppRepoWrite           repo.ClientAppRepository
	AuthQueryOAuth               query.AuthQueryOA
	AuthQueryDB                  query.AuthQuery
	MemberRepoRead               memberRepo.MemberRepository
	MemberRepoWrite              memberRepo.MemberRepository
	MemberQueryRead              memberQuery.MemberQuery
	MemberQueryWrite             memberQuery.MemberQuery
	CorporateContactQueryRead    corporateQuery.ContactQuery
	CorporateAccContactQueryRead corporateQuery.AccountContactQuery
	RefreshTokenRepo             repo.RefreshTokenRepository
	LoginAttemptRepo             repo.AttemptRepository
	LoginSessionRepo             repo.LoginSessionRepository
	AccessTokenGenerator         token.AccessTokenGenerator
	Hash                         memberModel.PasswordHasher
	RefreshTokenAge              string
	SpecialRefreshTokenAge       string
	EmailSpecialTokenAge         string
	LoginAttemptAge              string

	//Messaging
	QPublisher        service.QPublisher
	Topic             string
	IsProductionStage bool

	//Services
	AuthServices authServices.LDAPService

	SessionInfoRepo          sessionInfoRepo.SessionInfoRepository
	MerchantRepoRead         merchantRepoRead.MerchantRepository
	MerchantEmployeeRepoRead merchantRepoRead.MerchantEmployeeRepository

	GoogleVerifyCaptchaURL *url.URL
	StaticService          service.StaticServices
	ActivityService        service.ActivityServices
	B2cCFUrl               string
	NotificationService    service.NotificationServices
}

// NewAuthUseCase function for initialise auth use case implmentation model
func NewAuthUseCase(repository localConfig.ServiceRepository,
	queryParam localConfig.ServiceQuery,
	services localConfig.ServiceShared,
	params localConfig.AuthParameters,
	authQueryOAuth query.AuthQueryOA,
	authQueryDB query.AuthQuery) AuthUseCase {
	return &AuthUseCaseImpl{
		AuthQueryOAuth:               authQueryOAuth,
		AuthQueryDB:                  authQueryDB,
		ClientAppRepoRead:            repository.ClientAppRepoRead,
		ClientAppRepoWrite:           repository.ClientAppRepoWrite,
		MemberRepoRead:               repository.MemberRepository,
		MemberRepoWrite:              repository.MemberRepository,
		RefreshTokenRepo:             repository.RefreshTokenRepository,
		LoginAttemptRepo:             repository.AttemptRepositoryRedis,
		LoginSessionRepo:             repository.LoginSessionRepositoryRedis,
		SessionInfoRepo:              repository.SessionInfoRepo,
		MerchantRepoRead:             repository.MerchantRepository,
		MerchantEmployeeRepoRead:     repository.MerchantEmployeeRepository,
		MemberQueryRead:              queryParam.MemberQueryRead,
		MemberQueryWrite:             queryParam.MemberQueryWrite,
		CorporateContactQueryRead:    queryParam.CorporateContactQueryRead,
		CorporateAccContactQueryRead: queryParam.CorporateAccContactQueryRead,
		QPublisher:                   services.QPublisher,
		StaticService:                services.StaticService,
		ActivityService:              services.ActivityService,
		Hash:                         params.Hash,
		AuthServices:                 params.AuthServices,
		GoogleVerifyCaptchaURL:       params.GoogleVerifyCaptchaURL,
		AccessTokenGenerator:         params.AccessTokenGenerator,
		RefreshTokenAge:              params.RefreshTokenAge,
		B2cCFUrl:                     params.B2cCFUrl,
		LoginAttemptAge:              params.LoginAttemptAge,
		Topic:                        params.Topic,
		IsProductionStage:            params.IsProductionStage,
		NotificationService:          services.NotificationService,
		SpecialRefreshTokenAge:       params.SpecialRefreshTokenAge,
		EmailSpecialTokenAge:         params.EmailSpecialTokenAge,
	}
}
