package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/middleware"
	applicationsDelivery "github.com/Bhinneka/user-service/src/applications/v1/delivery"
	applicationsDeliveryV2 "github.com/Bhinneka/user-service/src/applications/v2/delivery"
	authDelivery "github.com/Bhinneka/user-service/src/auth/v1/delivery"
	authDeliveryV2 "github.com/Bhinneka/user-service/src/auth/v2/delivery"
	authDeliveryV3 "github.com/Bhinneka/user-service/src/auth/v3/delivery"
	clientDeliveryV1 "github.com/Bhinneka/user-service/src/client/v1/delivery"
	clientDeliveryV2 "github.com/Bhinneka/user-service/src/client/v2/delivery"
	corporateDeliveryV2 "github.com/Bhinneka/user-service/src/corporate/v2/delivery"
	documentDeliveryV2 "github.com/Bhinneka/user-service/src/document/v2/delivery"
	healthDelivery "github.com/Bhinneka/user-service/src/health/delivery"
	memberDelivery "github.com/Bhinneka/user-service/src/member/v1/delivery"
	memberDeliveryV2 "github.com/Bhinneka/user-service/src/member/v2/delivery"
	memberDeliveryV3 "github.com/Bhinneka/user-service/src/member/v3/delivery"
	memberDeliveryV4 "github.com/Bhinneka/user-service/src/member/v4/delivery"
	merchantDeliveryV2 "github.com/Bhinneka/user-service/src/merchant/v2/delivery"
	paymentsDelivery "github.com/Bhinneka/user-service/src/payments/v1/delivery"
	phoneAreaDelivery "github.com/Bhinneka/user-service/src/phone_area/v1/delivery"
	phoneAreaDeliveryV2 "github.com/Bhinneka/user-service/src/phone_area/v2/delivery"

	clientQuery "github.com/Bhinneka/user-service/src/client/v1/query"
	logDelivery "github.com/Bhinneka/user-service/src/log/v1/delivery"
	"github.com/Bhinneka/user-service/src/service"
	sessionInfoDelivery "github.com/Bhinneka/user-service/src/session/v1/delivery"
	shippingAddressDeliveryV2 "github.com/Bhinneka/user-service/src/shipping_address/v2/delivery"
	"github.com/labstack/echo"
	mid "github.com/labstack/echo/middleware"
)

const (
	//DefaultPort default http port
	DefaultPort = 8080
)

// HTTPServerMain main function for serving services over http
func (s *AppService) HTTPServerMain(publicKey *rsa.PublicKey, cq clientQuery.ClientQuery) {
	// construct Echo
	e := echo.New()
	e.Use(middleware.ServerHeader, middleware.Logger)

	serviceName := strings.TrimSuffix("user-service-"+os.Getenv("ENV"), "-PROD")
	if err := tracer.InitOpenTracing(os.Getenv("JAEGER_HOST"), serviceName); err == nil {
		e.Use(echo.WrapMiddleware(tracer.Middleware))
	}

	//e.Use(mid.Recover())

	if os.Getenv("ENABLE_CORS") == "true" {
		e.Use(mid.CORS())
	} else {
		e.Use(mid.CORSWithConfig(mid.CORSConfig{
			AllowOrigins: strings.Split(os.Getenv("CORS_DOMAIN"), ","),
		}))
	}

	if os.Getenv("DEVELOPMENT") == "1" {
		e.Debug = true
	}

	redisDB, ok := os.LookupEnv("REDIS_DB")
	if !ok {
		redisDB = "0"
	}

	redisConnection, err := redis.ConnectRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_TLS"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_PORT"), redisDB)
	if err != nil {
		log.Fatalf("redis error : %s", err.Error())
	}

	basicAuthUsername := os.Getenv("BASIC_USERNAME")
	basicAuthPassword := os.Getenv("BASIC_PASSWORD")
	uriV2 := "/api/v2"
	uriV3 := "/api/v3"
	uriV1 := "/api"

	basicAuthConfig := middleware.NewConfig(basicAuthUsername, basicAuthPassword)

	authHandler := authDelivery.NewHTTPHandler(s.AuthUseCase)
	authGroup := e.Group("/api/auth")
	authHandler.Mount(authGroup)

	// health endpoints
	healthHandler := healthDelivery.NewHTTPHandler(s.HealthUseCase)
	healthGroup := e.Group("/health")
	healthHandler.Mount(healthGroup)

	// member module handle
	memberHandler := memberDelivery.NewHTTPHandler(s.MemberUseCase)

	// member no token endpoints
	memberNoToken := e.Group(uriV1)
	memberHandler.Mount(memberNoToken)

	// member endpoints
	meGroup := e.Group("/api/me")
	meGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	memberHandler.MountMe(meGroup)

	// member endpoints /member
	memberGroup := e.Group("/api/member")
	memberGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	memberHandler.MountMember(memberGroup)

	memberGroupAuthorized := e.Group(uriV1)
	memberGroupAuthorized.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	memberHandler.MountAdmin(memberGroupAuthorized)

	// auth v2 endpoints
	googleAuthRedirectURL := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
	authHandlerV2 := authDeliveryV2.NewHTTPHandler(s.AuthUseCase, googleAuthRedirectURL)
	authGroupV2 := e.Group("/api/v2/auth")
	authHandlerV2.Mount(authGroupV2)

	authHandlerAdminV2 := authDeliveryV2.NewHTTPHandler(s.AuthUseCase, googleAuthRedirectURL)
	authGroupV2Admin := e.Group(uriV2)
	authGroupV2Admin.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	authHandlerAdminV2.MountAdmin(authGroupV2Admin)

	authHandlerV3 := authDeliveryV3.NewHTTPHandler(s.AuthUseCase, googleAuthRedirectURL)
	authGroupV3 := e.Group("/api/v3/auth")
	authHandlerV3.MountRoute(authGroupV3)

	// member v2 module handle
	memberHandlerV2 := memberDeliveryV2.NewHTTPHandler(s.MemberUseCase)

	// member no token v2 endpoints
	memberNoTokenV2 := e.Group(uriV2)
	memberHandlerV2.Mount(memberNoTokenV2)

	memberHandlerV3 := memberDeliveryV3.NewHTTPHandlerV3(s.MemberUseCase, s.AuthUseCase)
	memberNoTokenV3 := e.Group(uriV3)
	memberHandlerV3.Mount(memberNoTokenV3)

	memberSendirdGroupV3 := e.Group("/api/v3/sendbird")
	memberSendirdGroupV3.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	memberHandlerV3.MountSendbird(memberSendirdGroupV3)

	memberHandlerV4 := memberDeliveryV4.NewHTTPHandlerV4(s.MemberUseCase, s.AuthUseCase)
	memberSendbirdGroupV4 := e.Group("/api/v4/sendbird")
	memberSendbirdGroupV4.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	memberHandlerV4.MountSendbird(memberSendbirdGroupV4)

	// member v2 endpoints /me
	meGroup2 := e.Group("/api/v2/me")
	meGroup2.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	memberHandlerV2.MountMe(meGroup2)

	meGroup3 := e.Group("/api/v3/me")
	meGroup3.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	memberHandlerV3.MountMeV3(meGroup3)

	// member v2 endpoints /member
	memberGroupV2 := e.Group("/api/v2/member")
	memberGroupV2.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	memberHandlerV2.MountMember(memberGroupV2)

	// merchant v2 endpoints
	merchantHandlerV2 := merchantDeliveryV2.NewHTTPHandler(s.MerchantUseCase, s.MerchantAddressUseCase)
	merchantCMS := e.Group("/api/v2/merchant")
	merchantCMS.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	merchantHandlerV2.MountCMS(merchantCMS)

	merchantGroup := e.Group("/api/v2/merchant/me")
	merchantGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	merchantHandlerV2.MountMe(merchantGroup)

	merchantGroup2 := e.Group(uriV2)
	merchantGroup2.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	merchantHandlerV2.MountMerchant(merchantGroup2)

	merchantPublicV2 := e.Group(uriV2)
	merchantPublicV2.Use(middleware.BasicAuthWithConfig(cq))
	merchantHandlerV2.MountMerchantPublic(merchantPublicV2)

	// shippingAddress v2 endpoints
	shippingAddressHandlerV2 := shippingAddressDeliveryV2.NewHTTPHandler(s.ShippingAddressUseCase)
	shippingAddressGroup := e.Group("/api/v2/shipping-address")
	shippingAddressGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	shippingAddressHandlerV2.MountMe(shippingAddressGroup)

	shippingAddressGroup2 := e.Group("/api/v2/shipping-address")
	shippingAddressGroup2.Use(middleware.BearerVerify(publicKey, redisConnection, false, false))
	shippingAddressHandlerV2.MountShippingAddress(shippingAddressGroup2)

	// document v2 endpoints
	documentHandlerV2 := documentDeliveryV2.NewHTTPHandler(s.DocumentUseCase)
	documentGroup := e.Group("/api/v2/document")
	documentGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, false))
	documentHandlerV2.MountMe(documentGroup)

	documentTypeGroup := e.Group("/api/v2/document-type")
	documentTypeGroup.Use(middleware.BearerVerify(publicKey, redisConnection, false, true))
	documentHandlerV2.MountDocumentType(documentTypeGroup)

	// admin endpoints
	memberGroupAuthorizedV2 := e.Group(uriV2)
	memberGroupAuthorizedV2.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	memberHandlerV2.MountAdmin(memberGroupAuthorizedV2)

	// phone area endpoints
	phoneAreaHandler := phoneAreaDelivery.NewHTTPHandler(s.PhoneAreaUseCase)
	phoneAreaGroup := e.Group("/api/phone-area")
	phoneAreaGroup.Use(middleware.BasicAuth(basicAuthConfig))
	phoneAreaHandler.MountPhoneArea(phoneAreaGroup)

	// phone area v2 endpoints
	phoneAreaHandlerV2 := phoneAreaDeliveryV2.NewHTTPHandler(s.PhoneAreaUseCase)
	phoneAreaGroupV2 := e.Group("/api/v2/phone-area")
	phoneAreaGroupV2.Use(middleware.BasicAuthWithConfig(cq))
	phoneAreaHandlerV2.MountPhoneArea(phoneAreaGroupV2)

	// session info v1 endpoints
	sessionInfoHanlder := sessionInfoDelivery.NewHTTPHandler(s.SessionInfoUseCase)
	sessionInfoGroup := e.Group("/api/v1/session")
	sessionInfoGroup.Use(middleware.BasicAuth(basicAuthConfig))
	sessionInfoHanlder.MountInfo(sessionInfoGroup)

	// applications v1 endpoints
	applicationsHandler := applicationsDelivery.NewHTTPHandler(s.ApplicationsUseCase)
	applicationsGroup := e.Group("/api/v1/applications")
	applicationsGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	applicationsHandler.MountInfo(applicationsGroup)

	// applications v2 endpoints
	applicationsHandlerV2 := applicationsDeliveryV2.NewHTTPHandler(s.ApplicationsUseCase)
	applicationsGroupV2 := e.Group("/api/v2/applications")
	applicationsGroupV2.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	applicationsHandlerV2.MountInfo(applicationsGroupV2)

	// payment v1 endpoints
	paymentsHandler := paymentsDelivery.NewHTTPHandler(s.PaymentsUseCase)
	paymentsGroup := e.Group("/api/v1/payments")
	//paymentsGroup.Use(middleware.BasicAuth(basicAuthConfig))
	paymentsHandler.MountInfo(paymentsGroup)

	// corporate v2 endpoints
	corporateHandlerV2 := corporateDeliveryV2.NewHTTPHandler(s.CorporateUseCase)
	corporateGroup2 := e.Group("/api/v2/corporate")
	corporateGroup2.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	corporateHandlerV2.MountCorporate(corporateGroup2)
	activityService := service.NewActivityService("v2")
	// client endpoint

	clientHandlerV1 := clientDeliveryV1.NewHTTPHandler(s.AuthUseCase, activityService, s.ClientUseCase)
	clientGroup := e.Group("/v1/client")
	clientHandlerV1.Mount(clientGroup, middleware.BasicAuth(basicAuthConfig))

	clientHandlerV2 := clientDeliveryV2.NewHTTPHandler(s.AuthUseCase, activityService, s.ClientUseCase, s.ClientV2UseCase)
	clientV2Group := e.Group("/v2/client")
	clientHandlerV2.Mount(clientV2Group, middleware.BasicAuth(basicAuthConfig))

	logHandlerV1 := logDelivery.NewHTTPHandler(activityService, s.LogUseCase)
	logGroup := e.Group("v1/log")
	logGroup.Use(middleware.BearerVerify(publicKey, redisConnection, true, true))
	logHandlerV1.Mount(logGroup)

	// set REST port
	var port uint16
	if portEnv, ok := os.LookupEnv("PORT"); ok {
		portInt, err := strconv.Atoi(portEnv)
		if err != nil {
			port = DefaultPort
		} else {
			port = uint16(portInt)
		}
	} else {
		port = DefaultPort
	}

	listenerPort := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(listenerPort))
}
