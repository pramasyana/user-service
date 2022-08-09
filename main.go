package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/config/rsa"
	"github.com/Bhinneka/user-service/helper"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	authToken "github.com/Bhinneka/user-service/src/auth/v1/token"
	clientQuery "github.com/Bhinneka/user-service/src/client/v1/query"
	corporateRepo "github.com/Bhinneka/user-service/src/corporate/v2/repo"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	paymentRepo "github.com/Bhinneka/user-service/src/payments/v1/repo"
	"github.com/Bhinneka/user-service/src/service"
	"github.com/Bhinneka/user-service/src/shared"
	sharedRepository "github.com/Bhinneka/user-service/src/shared/repository"
	shippingAddressRepo "github.com/Bhinneka/user-service/src/shipping_address/v2/repo"
	"github.com/getsentry/raven-go"
	config "github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {

	ctx := "user_service_main"

	err := config.Load(".env")
	if err != nil {
		fmt.Println(".env is not loaded properly")
		os.Exit(2)
	}

	// set DSN here
	isSentry, _ := strconv.ParseBool(os.Getenv("SENTRY"))

	// check setry to send or not
	if isSentry {
		dsn, ok := os.LookupEnv("SENTRY_DSN")
		if !ok {
			fmt.Println(".env is not loaded properly")
			os.Exit(2)
		}

		raven.SetDSN(dsn)
	}

	// redis env
	redisHost, ok := os.LookupEnv("REDIS_HOST")
	if !ok {
		fmt.Println("redis host is not loaded")
		os.Exit(2)
	}

	redisPort, ok := os.LookupEnv("REDIS_PORT")
	if !ok {
		fmt.Println("redis port is not loaded")
		os.Exit(2)
	}

	redisAuth, ok := os.LookupEnv("REDIS_PASSWORD")
	if !ok {
		fmt.Println("redis password is not loaded")
		os.Exit(2)
	}

	redisTLS, ok := os.LookupEnv("REDIS_TLS")
	if !ok {
		fmt.Println("redis TLS is not loaded")
		os.Exit(2)
	}

	redisDB := "0"
	redisDB, _ = os.LookupEnv("REDIS_DB")
	// redis env

	publicKey, err := rsa.InitPublicKey()
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "private_key")
		os.Exit(1)
	}

	// dolphin consumer status enable env
	enableConsumerDolphin, ok := os.LookupEnv("ENABLE_CONSUMER_MEMBER_DOLPHIN")
	if !ok {
		fmt.Println("consumer for dolphin member status is not loaded")
		os.Exit(2)
	}

	// initiate database and other connections
	readDB := localConfig.ReadPostgresDB()
	writeDB := localConfig.WritePostgresDB()

	//dolphin service
	dolphinService, err := service.NewDolphinService()
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "construct_dolphin_service")
		os.Exit(1)
	}

	redisConnection, err := redis.ConnectRedis(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_TLS"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_PORT"), redisDB)
	if err != nil {
		helper.Log(log.ErrorLevel, err.Error(), ctx, "redis_connection")
		os.Exit(1)
	}

	privateKey, _ := rsa.InitPrivateKey()
	tokenAge := golib.GetEnvDurationOrFail(ctx, "parse_token_age", "ACCESS_TOKEN_AGE")

	refreshTokenAge := golib.GetEnvDurationOrFail(ctx, "parse_refresh_token_age", "REFRESH_TOKEN_AGE")

	specialTokenAge := golib.GetEnvDurationOrFail(ctx, "parse_token_age", "SPECIAL_ACCESS_TOKEN_AGE")

	specialRefreshTokenAge := golib.GetEnvDurationOrFail(ctx, "parse_refresh_token_age", "SPECIAL_REFRESH_TOKEN_AGE")

	emailForSpecialToken := golib.GetEnvOrFail(ctx, "find_azure_config_login_url", "EMAIL_SPECIAL_TOKEN_AGE")

	// define parent repository from shared
	sRepository := sharedRepository.NewRepository(readDB, writeDB)

	// initial repository dolphin log
	dolphinLogRepository := memberRepo.NewDolphinLogRepoPostgres(writeDB)

	// initial member repository

	// initial member query
	memberQueryWrite := memberQuery.NewMemberQueryPostgres(writeDB)

	// initial member query
	memberRepoWrite := memberRepo.NewMemberRepoPostgres(sRepository)

	paymentRepo := paymentRepo.NewPaymentsRepoPostgres(sRepository)

	// initial repository account
	accountRepository := corporateRepo.NewAccountRepoPostgres(sRepository)

	// initital repository account temporary
	accountTemporaryRepository := corporateRepo.NewAccountTemporaryRepoPostgres(sRepository)

	// initital repository account contact
	accountContactRepository := corporateRepo.NewAccountContactRepoPostgres(sRepository)

	// initital repository contact
	contactRepository := corporateRepo.NewContactRepoPostgres(sRepository)

	// initital repository address
	addressRepository := corporateRepo.NewAddressRepoPostgres(sRepository)

	// initital repository phone
	phoneRepository := corporateRepo.NewPhoneRepoPostgres(sRepository)

	// initital repository document
	documentRepository := corporateRepo.NewDocumentRepoPostgres(sRepository)

	// initital repository contact address
	contactAddressRepository := corporateRepo.NewContactAddressRepoPostgres(sRepository)

	// initital repository contact npwp
	contactNpwpRepository := corporateRepo.NewContactNpwpRepoPostgres(sRepository)

	// initital repository contact temp
	contactTempRepository := corporateRepo.NewContactTempRepoPostgres(sRepository)

	// initital repository leads
	leadsRepository := corporateRepo.NewLeadsRepoPostgres(sRepository)

	// initital repository merchant
	merchantRepository := merchantRepo.NewMerchantRepoPostgres(sRepository)

	// initital repository merchant doc
	merchantDocumentRepository := merchantRepo.NewMerchantDocumentRepoPostgres(sRepository)

	// initital repository merchant bank
	merchantBankRepository := merchantRepo.NewMerchantBankRepoPostgres(sRepository)

	// initital repository merchant employee
	merchantEmployeeRepository := merchantRepo.NewMerchantEmployeeRepoPostgres(sRepository)

	// initital repository shipping address
	shippingAddressRepository := shippingAddressRepo.NewShippingAddressRepoPostgres(sRepository)

	contactDocumentRepository := corporateRepo.NewContactDocumentRepoPostgres(sRepository)

	// initital repository session redis repo
	loginSessionRedisRepo := authRepo.NewLoginSessionRepositoryRedis(redisConnection)

	// initital jwt generator
	accessTokenGenerator := authToken.NewJwtGenerator(privateKey, tokenAge, refreshTokenAge, specialTokenAge, specialRefreshTokenAge, loginSessionRedisRepo, emailForSpecialToken)

	kafkaBroker1 := golib.GetEnvOrFail(ctx, "kafka_find_host_1", "KAFKA_BHINNEKA_BROKER_1")
	kafkaBroker2 := golib.GetEnvOrFail(ctx, "kafka_find_host_2", "KAFKA_BHINNEKA_BROKER_2")
	kafkaBroker3 := golib.GetEnvOrFail(ctx, "kafka_find_host_3", "KAFKA_BHINNEKA_BROKER_3")
	brokers := []string{kafkaBroker1, kafkaBroker2, kafkaBroker3}
	kafkaMessaging, _ := service.NewKafkaPublisher(kafkaBroker1, kafkaBroker2, kafkaBroker3)

	app := MakeHandler(readDB, writeDB, kafkaMessaging)
	clientAppQuery := clientQuery.NewClientAppQuery(readDB)

	// redis subscriber
	redisPubSubConfig := &shared.RedisPubSubConfig{
		Host:     redisHost,
		Password: redisAuth,
		Port:     redisPort,
		UseTLS:   redisTLS,
		UseDB:    redisDB,
	}

	redisPubSubClient := shared.NewRedisPubSub(redisPubSubConfig)

	pathSchema := "schema/json"
	jsonschema.Load(pathSchema)
	serviceRepo := localConfig.ServiceRepository{
		CorporateAccountRepository:         accountRepository,
		CorporateAccountTempRepository:     accountTemporaryRepository,
		CorporateAccountContactRepository:  accountContactRepository,
		CorporateContactRepository:         contactRepository,
		CorporateAddressRepository:         addressRepository,
		CorporatePhoneRepository:           phoneRepository,
		CorporateDocumentRepository:        documentRepository,
		CorporateContactNPWPRepository:     contactNpwpRepository,
		CorporateContactAddressRepository:  contactAddressRepository,
		CorporateContactTempRepository:     contactTempRepository,
		CorporateLeadsRepository:           leadsRepository,
		MerchantRepository:                 merchantRepository,
		MerchantDocumentRepository:         merchantDocumentRepository,
		MerchantBankRepository:             merchantBankRepository,
		MerchantEmployeeRepository:         merchantEmployeeRepository,
		ShippingAddressRepository:          shippingAddressRepository,
		CorporateContactDocumentRepository: contactDocumentRepository,
		MemberRepository:                   memberRepoWrite,
		PaymentsRepository:                 paymentRepo,
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.HTTPServerMain(publicKey, clientAppQuery)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		app.GRPCServerMain(publicKey, kafkaMessaging)
	}()

	if enableConsumerDolphin == "true" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumeKafka(dolphinLogRepository, dolphinService, brokers)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		dispatchWorker(app, brokers)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeKafkaGWS(sRepository, merchantRepository, merchantDocumentRepository, merchantBankRepository, accessTokenGenerator, kafkaMessaging, brokers)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumeKafkaShark(serviceRepo, brokers)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		redisSubscribe(redisPubSubClient, memberQueryWrite)
	}()

	// Wait All services to end
	wg.Wait()
}
