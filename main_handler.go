package main

import (
	"crypto/sha1"
	"database/sql"
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/Bhinneka/golib"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/config/rsa"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/service"
	log "github.com/sirupsen/logrus"

	healthQuery "github.com/Bhinneka/user-service/src/health/query"
	healthUseCase "github.com/Bhinneka/user-service/src/health/usecase"

	authQuery "github.com/Bhinneka/user-service/src/auth/v1/query"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	authToken "github.com/Bhinneka/user-service/src/auth/v1/token"
	authUseCase "github.com/Bhinneka/user-service/src/auth/v1/usecase"

	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	memberUseCase "github.com/Bhinneka/user-service/src/member/v1/usecase"

	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	corporateRepository "github.com/Bhinneka/user-service/src/corporate/v2/repo"
	corporateUseCase "github.com/Bhinneka/user-service/src/corporate/v2/usecase"

	phoneAreaQuery "github.com/Bhinneka/user-service/src/phone_area/v1/query"
	phoneAreaUseCase "github.com/Bhinneka/user-service/src/phone_area/v1/usecase"

	sessionInfoQuery "github.com/Bhinneka/user-service/src/session/v1/query"
	sessionInfoRepository "github.com/Bhinneka/user-service/src/session/v1/repo"
	sessionInfoUseCase "github.com/Bhinneka/user-service/src/session/v1/usecase"

	applicationRepository "github.com/Bhinneka/user-service/src/applications/v1/repo"
	applicationsUseCase "github.com/Bhinneka/user-service/src/applications/v1/usecase"
	authServices "github.com/Bhinneka/user-service/src/auth/v1/service"

	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	merchantUseCase "github.com/Bhinneka/user-service/src/merchant/v2/usecase"

	sharedRepository "github.com/Bhinneka/user-service/src/shared/repository"

	shippingAddressRepository "github.com/Bhinneka/user-service/src/shipping_address/v2/repo"
	shippingAddressUseCase "github.com/Bhinneka/user-service/src/shipping_address/v2/usecase"

	documentRepository "github.com/Bhinneka/user-service/src/document/v2/repo"
	documentUseCase "github.com/Bhinneka/user-service/src/document/v2/usecase"

	clientUseCase "github.com/Bhinneka/user-service/src/client/v1/usecase"
	clientV2UseCase "github.com/Bhinneka/user-service/src/client/v2/usecase"
	logUseCase "github.com/Bhinneka/user-service/src/log/v1/usecase"

	paymentRepo "github.com/Bhinneka/user-service/src/payments/v1/repo"
	paymentUseCase "github.com/Bhinneka/user-service/src/payments/v1/usecase"
)

const (
	parseRefreshTokenAgeText = "parse_refresh_token_age"
)

// AppService data structure
type AppService struct {
	HealthUseCase          healthUseCase.HealthUseCase
	AuthUseCase            authUseCase.AuthUseCase
	MemberUseCase          memberUseCase.MemberUseCase
	PhoneAreaUseCase       phoneAreaUseCase.PhoneAreaUseCase
	SessionInfoUseCase     sessionInfoUseCase.SessionInfoUseCase
	ApplicationsUseCase    applicationsUseCase.ApplicationsUseCase
	MerchantUseCase        merchantUseCase.MerchantUseCase
	ShippingAddressUseCase shippingAddressUseCase.ShippingAddressUseCase
	DocumentUseCase        documentUseCase.DocumentUseCase
	CorporateUseCase       corporateUseCase.CorporateUseCase
	MerchantAddressUseCase merchantUseCase.MerchantAddressUseCase
	ClientUseCase          clientUseCase.ClientUsecase
	ClientV2UseCase        clientV2UseCase.ClientUsecase
	LogUseCase             logUseCase.LogUsecase
	PaymentsUseCase        paymentUseCase.PaymentsUseCase
}

//MakeHandler function, Service's Constructor
func MakeHandler(readDB, writeDB *sql.DB, kafkaMessaging service.QPublisher) *AppService {
	ctx := "make_handler"

	privateKey, err := rsa.InitPrivateKey()

	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "private_key")
		os.Exit(1)
	}

	redisDB := "0"
	redisDB, _ = os.LookupEnv("REDIS_DB")

	redisConnection, err := redis.ConnectRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_TLS"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_PORT"), redisDB)
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "redis_connection")
		os.Exit(1)
	}

	tokenAge := golib.GetEnvDurationOrFail(ctx, "parse_token_age", "ACCESS_TOKEN_AGE")

	refreshTokenAge := golib.GetEnvDurationOrFail(ctx, parseRefreshTokenAgeText, "REFRESH_TOKEN_AGE")

	specialTokenAge := golib.GetEnvDurationOrFail(ctx, "parse_token_age", "SPECIAL_ACCESS_TOKEN_AGE")

	specialRefreshTokenAge := golib.GetEnvDurationOrFail(ctx, parseRefreshTokenAgeText, "SPECIAL_REFRESH_TOKEN_AGE")

	specialRefreshTokenAgeString := golib.GetEnvOrFail(ctx, parseRefreshTokenAgeText, "SPECIAL_REFRESH_TOKEN_AGE")

	emailForSpecialToken := golib.GetEnvOrFail(ctx, "find_azure_config_login_url", "EMAIL_SPECIAL_TOKEN_AGE")

	refreshTokenAgeStr := golib.GetEnvOrFail(ctx, parseRefreshTokenAgeText, "REFRESH_TOKEN_AGE")

	loginAttemptAge := golib.GetEnvOrFail(ctx, "parse_login_attempt_age", "LOGIN_ATTEMPT_AGE")

	tokenActivationAge := golib.GetEnvDurationOrFail(ctx, "parse_token_expiration_age", "TOKEN_ACTIVATION_AGE")

	// messaging system
	// NSQ
	NSQServer := golib.GetEnvOrFail(ctx, "find_nsq_server_config", "NSQ_LOOKUP")

	topic := golib.GetEnvOrFail(ctx, "find_nsq_topic_config", "NSQ_TOPIC_NOTIFICATION")

	kafkaUserServiceTopic := golib.GetEnvOrFail(ctx, helper.TextFindServerKafkaConfig, "KAFKA_USER_SERVICE_TOPIC")

	// set password hash
	passwordHasher := memberModel.NewPBKDF2Hasher(memberModel.SaltSize, memberModel.SaltSize, memberModel.IterationsCount, sha1.New)

	// set OAuth social media url
	azureLoginBaseURL := golib.GetEnvOrFail(ctx, "find_azure_config_login_url", "AD_LOGIN_URL")

	azureLoginParsedURL, err := url.Parse(azureLoginBaseURL)
	if err != nil {
		err := errors.New("invalid AD_LOGIN_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_azure_config_login_url")
		os.Exit(1)
	}

	azureBaseURL := golib.GetEnvOrFail(ctx, "find_azure_config_resource", "AD_RESOURCE")
	azureParsedURL, err := url.Parse(azureBaseURL)
	if err != nil {
		err := errors.New("invalid AD_LOGIN_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_azure_config_resource")
		os.Exit(1)
	}

	facebookClientID := golib.GetEnvOrFail(ctx, "find_facebook_config_client_id", "FACEBOOK_CLIENT_ID")

	facebookClientSecret := golib.GetEnvOrFail(ctx, "find_facebook_config_client_secret", "FACEBOOK_CLIENT_SECRET")

	facebookBaseURL := golib.GetEnvOrFail(ctx, "find_facebook_config_login_url", "FACEBOOK_LOGIN_URL")

	facebookParsedURL, err := url.Parse(facebookBaseURL)
	if err != nil {
		err := errors.New("invalid FACEBOOK_LOGIN_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_facebook_config")
		os.Exit(1)
	}

	googleBaseURL := golib.GetEnvOrFail(ctx, "find_google_config_env", "GOOGLE_LOGIN_URL")

	googleParsedURL, err := url.Parse(googleBaseURL)
	if err != nil {
		err := errors.New("invalid GOOGLE_LOGIN_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_google_config")
		os.Exit(1)
	}

	googleTokenBaseURL := golib.GetEnvOrFail(ctx, "find_google_config_auth_env", "GOOGLE_AUTH_URL")

	googleTokenParsedURL, err := url.Parse(googleTokenBaseURL)
	if err != nil {
		err := errors.New("invalid GOOGLE_AUTH_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_google_config")
		os.Exit(1)
	}

	googleOAuthBaseURL := golib.GetEnvOrFail(ctx, "find_google_oauth_config_auth_env", "GOOGLE_OAUTH_URL")

	googleOauthParsedURL, err := url.Parse(googleOAuthBaseURL)
	if err != nil {
		err := errors.New("invalid GOOGLE_OAUTH_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_google_config")
		os.Exit(1)
	}

	appleTokenBaseURL := golib.GetEnvOrFail(ctx, "find_apple_config_auth_env", "APPLE_AUTH_URL")

	appleTokenParsedURL, err := url.Parse(appleTokenBaseURL)
	if err != nil {
		err := errors.New("invalid GOOGLE_AUTH_URL in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_apple_config")
		os.Exit(1)
	}

	googleClientID := golib.GetEnvOrFail(ctx, "find_google_config_client_id_env", "GOOGLE_CLIENT_ID")

	googleClientSecret := golib.GetEnvOrFail(ctx, "find_google_config", "GOOGLE_CLIENT_SECRET")

	azureADTenanID := golib.GetEnvOrFail(ctx, "find_azure_config_tenant_id", "AD_TENANT_ID")

	azureADClientID := golib.GetEnvOrFail(ctx, "find_azure_config_client_id", "AD_CLIENT_ID")

	azureADClientSecret := golib.GetEnvOrFail(ctx, "find_azure_config_client_secret", "AD_CLIENT_SECRET")

	azureADResource := golib.GetEnvOrFail(ctx, "find_azure_config_ad_resource", "AD_RESOURCE")

	appleTeamID := golib.GetEnvOrFail(ctx, "find_apple_config_team_id_env", "APPLE_TEAM_ID")

	appleKeyID := golib.GetEnvOrFail(ctx, "find_apple_config_key_id_env", "APPLE_KEY_ID")

	isProductionStageString := golib.GetEnvOrFail(ctx, "find_stage_config_env", "IS_PRODUCTION_STAGE")

	isProductionStage, err := strconv.ParseBool(isProductionStageString)
	if err != nil {
		err := errors.New("IS_PRODUCTION_STAGE should be boolean")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "find_stage_config_parse")
		os.Exit(1)
	}

	sturgeonCFUrl := golib.GetEnvOrFail(ctx, "find_sturgeon_url_config", "STURGEON_CF_URL")

	b2cCFUrl := golib.GetEnvOrFail(ctx, "find_sturgeon_url_config", "B2C_CF_URL")

	googleVerifyCaptchaURL := golib.GetEnvOrFail(ctx, "find_google_config_env", "GOOGLE_VERIFY_CAPTCHA")

	googleVerifyCaptchaParsedURL, err := url.Parse(googleVerifyCaptchaURL)
	if err != nil {
		err := errors.New("invalid GOOGLE_VERIFY_CAPTCHA in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "parse_google_captcha_config")
		os.Exit(1)
	}

	resendActivationAttemptAge := golib.GetEnvOrFail(ctx, "parse_resend_activation_attempt_age", "RESEND_ACTIVATION_ATTEMPT_AGE")

	resendActivationAttemptAgeRequest := golib.GetEnvOrFail(ctx, "parse_resend_activation_attempt_age_request", "RESEND_ACTIVATION_ATTEMPT_AGE_REQUEST")

	//service
	//nsqDispatcher := service.NewNSQDispatcher()

	//activity service
	activityService := service.NewActivityService("v2")

	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "kafka_initialization")
		os.Exit(1)
	}

	//static service
	staticService, err := service.NewStaticService()
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "construct_static_service")
		os.Exit(1)
	}

	//merchant service
	merchantService, err := service.NewMerchantService(kafkaMessaging, activityService)
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "construct_merchant_service")
		os.Exit(1)
	}

	//barracuda service
	barracudaService, err := service.NewBarracudaService()
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "construct_barracuda_service")
		os.Exit(1)
	}

	//sendbird service
	sendbirdService, err := service.NewSendbirdService()
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "construct_sendbird_service")
		os.Exit(1)
	}

	//upload service
	uploadService, err := service.NewUploadService()
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "construct_upload_service")
		os.Exit(1)
	}

	// define parent repository from shared
	sRepository := sharedRepository.NewRepository(readDB, writeDB)

	// connection initializing
	hQuery := healthQuery.NewHealthQueryImpl("I'm fine!!", readDB, redisConnection, NSQServer, topic)
	aRepoRead := authRepo.NewClientAppRepoPostgres(writeDB)
	aRepoWrite := authRepo.NewClientAppRepoPostgres(writeDB)
	aQuery := authQuery.NewAuthQueryPostgres(writeDB)
	oAuthConfig := localConfig.OAuthService{
		AzureLoginBaseURL:    azureLoginParsedURL,
		AzureBaseURL:         azureParsedURL,
		FacebookBaseURL:      facebookParsedURL,
		FacebookTokenBaseURL: facebookParsedURL,
		GoogleBaseURL:        googleParsedURL,
		GoogleBaseTokenURL:   googleTokenParsedURL,
		GoogleOAuthBaseURL:   googleOauthParsedURL,
		AppleBaseURL:         appleTokenParsedURL,
		FBClientID:           facebookClientID,
		FBClientSecret:       facebookClientSecret,
		GoogleClientID:       googleClientID,
		GoogleClientSecret:   googleClientSecret,
		AzureADTenanID:       azureADTenanID,
		AzureADClientID:      azureADClientID,
		AzureADClientSecret:  azureADClientSecret,
		AzureADResource:      azureADResource,
		AppleTeamID:          appleTeamID,
		AppleKeyID:           appleKeyID,
	}

	aQueryOAuth := authQuery.NewAuthQueryOAuth(oAuthConfig)
	refreshTokenRepo := authRepo.NewRefreshTokenRepositoryRedis(redisConnection)
	tokenActivationRepo := memberRepo.NewTokenActivationRepoRedis(redisConnection)
	attemptRepo := authRepo.NewAttemptRepositoryRedis(redisConnection)
	mRepo := memberRepo.NewMemberRepoPostgres(sRepository)
	mMFARepo := memberRepo.NewMemberMFARepoPostgres(sRepository)
	mRepoRedis := memberRepo.NewMemberRepoRedis(redisConnection)
	mAdditionalRepo := memberRepo.NewMemberAdditionalInfoRepoPostgres(sRepository)
	mQueryRead := memberQuery.NewMemberQueryPostgres(readDB)
	mQueryWrite := memberQuery.NewMemberQueryPostgres(writeDB)
	mMFAQueryRead := memberQuery.NewMemberMFAQueryPostgres(readDB)

	cContactQueryRead := corporateQuery.NewContactQueryPostgres(readDB)
	cAccountContactQueryRead := corporateQuery.NewAccountContactQueryPostgres(readDB)
	contactRepo := corporateRepository.NewContactRepoPostgres(sRepository)

	aRepoSessionInfo := sessionInfoRepository.NewSessionInfoRepoPostgres(writeDB)
	aServiceLdap, _ := authServices.NewLDAPService(os.Getenv("LDAP_SERVER"),
		os.Getenv("LDAP_BIND"),
		os.Getenv("LDAP_PASSWORD"),
		os.Getenv("LDAP_BASE_DN"),
		os.Getenv("LDAP_FILTER_DN"),
		mAdditionalRepo, mQueryRead)
	sessionQuery := sessionInfoQuery.NewSessionInfoQueryPostgres(sRepository)
	loginSessionRedisRepo := authRepo.NewLoginSessionRepositoryRedis(redisConnection)
	pAreaQueryRead := phoneAreaQuery.NewPhoneAreaQueryPostgres(readDB)
	paymentRepo := paymentRepo.NewPaymentsRepoPostgres(sRepository)
	merchantRepository := merchantRepo.NewMerchantRepoPostgres(sRepository)
	merchantBankRepository := merchantRepo.NewMerchantBankRepoPostgres(sRepository)
	merchantEmployeeRepository := merchantRepo.NewMerchantEmployeeRepoPostgres(sRepository)
	merchantDocumentRepository := merchantRepo.NewMerchantDocumentRepoPostgres(sRepository)
	merchantAddressRepository := merchantRepo.NewMerchantAddressRepoPostgres(sRepository)
	shippingAddressRepo := shippingAddressRepository.NewShippingAddressRepoPostgres(sRepository)
	shippingAddressRedisRepo := shippingAddressRepository.NewShippingAddressRepoRedis(redisConnection)
	applicationRepo := applicationRepository.NewApplicationRepoPostgres(sRepository)
	documentRepo := documentRepository.NewDocumentRepoPostgres(sRepository, uploadService)
	documentTypeRepo := documentRepository.NewDocumentTypeRepoPostgres(sRepository)
	jwtGenerator := authToken.NewJwtGenerator(privateKey, tokenAge, refreshTokenAge, specialTokenAge, specialRefreshTokenAge, loginSessionRedisRepo, emailForSpecialToken)
	notificationService := service.NewNotificationService(jwtGenerator)

	serviceRepo := localConfig.ServiceRepository{
		ClientAppRepoRead:              aRepoRead,
		ClientAppRepoWrite:             aRepoWrite,
		DocumentRepository:             documentRepo,
		MerchantRepository:             merchantRepository,
		MerchantDocumentRepository:     merchantDocumentRepository,
		MerchantBankRepository:         merchantBankRepository,
		MerchantEmployeeRepository:     merchantEmployeeRepository,
		MerchantAddressRepository:      merchantAddressRepository,
		ShippingAddressRepository:      shippingAddressRepo,
		ShippingAddressRedisRepository: shippingAddressRedisRepo,
		MemberRepository:               mRepo,
		MemberMFARepository:            mMFARepo,
		MemberRedisRepository:          mRepoRedis,
		TokenActivationRepoRedis:       tokenActivationRepo,
		AttemptRepositoryRedis:         attemptRepo,
		LoginSessionRepositoryRedis:    loginSessionRedisRepo,
		Repository:                     sRepository,
		RefreshTokenRepository:         refreshTokenRepo,
		SessionInfoRepo:                aRepoSessionInfo,
		PaymentsRepository:             paymentRepo,
	}

	serviceQuery := localConfig.ServiceQuery{
		CorporateContactQueryRead:    cContactQueryRead,
		CorporateAccContactQueryRead: cAccountContactQueryRead,
		MemberQueryRead:              mQueryRead,
		MemberQueryWrite:             mQueryWrite,
		MemberMFAQueryRead:           mMFAQueryRead,
		SessionInfoQuery:             sessionQuery,
	}

	serviceShared := localConfig.ServiceShared{
		StaticService:       staticService,
		UploadService:       uploadService,
		MerchantService:     merchantService,
		ActivityService:     activityService,
		BarracudaService:    barracudaService,
		QPublisher:          kafkaMessaging,
		NotificationService: notificationService,
		SendbirdService:     sendbirdService,
	}
	// all usecase
	hUseCase := healthUseCase.NewHealthUseCase(hQuery)

	membershipParameters := localConfig.MembershipParameters{
		Hash:                              passwordHasher,
		TokenActivationExpiration:         tokenActivationAge,
		ResendActivationAttemptAge:        resendActivationAttemptAge,
		ResendActivationAttemptAgeRequest: resendActivationAttemptAgeRequest,
		Topic:                             kafkaUserServiceTopic,
		IsProductionStage:                 isProductionStage,
		SturgeonCFUrl:                     sturgeonCFUrl,
		B2cCFUrl:                          b2cCFUrl,
		AccessTokenGenerator:              jwtGenerator,
	}

	authParameters := localConfig.AuthParameters{
		Hash:                   passwordHasher,
		AuthServices:           aServiceLdap,
		GoogleVerifyCaptchaURL: googleVerifyCaptchaParsedURL,
		AccessTokenGenerator:   jwtGenerator,
		RefreshTokenAge:        refreshTokenAgeStr,
		B2cCFUrl:               b2cCFUrl,
		LoginAttemptAge:        loginAttemptAge,
		Topic:                  kafkaUserServiceTopic,
		IsProductionStage:      isProductionStage,
		EmailSpecialTokenAge:   emailForSpecialToken,
		SpecialRefreshTokenAge: specialRefreshTokenAgeString,
	}

	aUseCase := authUseCase.NewAuthUseCase(serviceRepo, serviceQuery, serviceShared, authParameters, aQueryOAuth, aQuery)
	pAreaUseCase := phoneAreaUseCase.NewPhoneAreaUseCase(pAreaQueryRead)
	sessionUseCase := sessionInfoUseCase.NewSessionInfoUseCase(sessionQuery, aRepoSessionInfo)
	appsUseCase := applicationsUseCase.NewApplicationsUseCase("APPLICATIONS_JSON", applicationRepo)
	merchantUseCases := merchantUseCase.NewMerchantUseCase(serviceRepo, serviceShared, jwtGenerator, serviceQuery)
	merchantAddressUseCase := merchantUseCase.NewMerchantAddressUseCase(serviceRepo, serviceShared)
	shippingAddressUseCase := shippingAddressUseCase.NewShippingAddressUseCase(serviceRepo, serviceShared)
	documentUseCase := documentUseCase.NewDocumentUseCase(documentRepo, documentTypeRepo, mRepo, sRepository, "DOCUMENTS_JSON")
	corporateUseCase := corporateUseCase.NewCorporateUseCase(contactRepo, cContactQueryRead, serviceShared)
	clientUseCase := clientUseCase.NewClientUsecase(loginSessionRedisRepo, refreshTokenRepo, mQueryRead)
	clientV2UseCase := clientV2UseCase.NewClientUsecase(loginSessionRedisRepo, refreshTokenRepo, mQueryRead, cContactQueryRead)
	logUsecase := logUseCase.NewLogUsecase(serviceShared)
	paymentUseCase := paymentUseCase.NewPaymentsUseCase(serviceRepo, serviceQuery)

	mUseCase := memberUseCase.NewMemberUseCase(serviceRepo, serviceQuery, serviceShared, membershipParameters, aUseCase)

	return &AppService{
		HealthUseCase:          hUseCase,
		MemberUseCase:          mUseCase,
		AuthUseCase:            aUseCase,
		PhoneAreaUseCase:       pAreaUseCase,
		SessionInfoUseCase:     sessionUseCase,
		ApplicationsUseCase:    appsUseCase,
		MerchantUseCase:        merchantUseCases,
		ShippingAddressUseCase: shippingAddressUseCase,
		DocumentUseCase:        documentUseCase,
		CorporateUseCase:       corporateUseCase,
		MerchantAddressUseCase: merchantAddressUseCase,
		ClientUseCase:          clientUseCase,
		ClientV2UseCase:        clientV2UseCase,
		LogUseCase:             logUsecase,
		PaymentsUseCase:        paymentUseCase,
	}
}
