package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	localConfig "github.com/Bhinneka/user-service/config"
	mockMemberRepo "github.com/Bhinneka/user-service/mocks/src/member/v1/repo"
	mockToken "github.com/Bhinneka/user-service/src/auth/v1/token/mocks"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/repo"
	mockMerchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo/mocks"
	serviceMock "github.com/Bhinneka/user-service/src/service/mocks"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	sharedMock "github.com/Bhinneka/user-service/src/shared/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

const (
	defEmail          = "pian.mutakin@bhinneka.com"
	defMerchantName   = "Nafara Collection"
	defUserID         = "USR1111"
	defMerchantGroup  = "perorangan"
	defPhoneNumber    = "081282928292"
	defDescription    = "some description"
	defCompanyName    = "my company"
	defAddress        = "my address"
	defPIC            = "my pic"
	defImage          = "https://someurl.com/image.png"
	defPicOccupation  = "manager"
	defName           = "some name"
	defNpwp           = "012909910102909"
	defIPAddr         = "127.0.0.1"
	defaultMerchantID = "MCH001"
	defDocumentID     = "DOC001"
)

var (
	errDefault      = fmt.Errorf("default error")
	defUserAttr     = &model.MerchantUserAttribute{UserID: defUserID, UserIP: defIPAddr}
	defEmailContent = serviceModel.Template{Content: "some email content"}
	defInput        = model.B2CMerchantCreateInput{
		MerchantEmail:         defEmail,
		MerchantName:          defMerchantName,
		UserID:                defUserID,
		GenderPicString:       "MALE",
		MerchantGroup:         model.MicroString,
		Documents:             []model.B2CMerchantDocumentInput{},
		PhoneNumber:           defPhoneNumber,
		BusinessType:          "perusahaan",
		MerchantDescription:   defDescription,
		CompanyName:           defCompanyName,
		MerchantAddress:       defAddress,
		Pic:                   defPIC,
		PicKtpFile:            defImage,
		PicOccupation:         defPicOccupation,
		MobilePhoneNumber:     defPhoneNumber,
		NpwpFile:              defImage,
		NpwpHolderName:        defName,
		Npwp:                  defNpwp,
		AccountNumber:         defPhoneNumber,
		AccountHolderName:     defName,
		BankBranch:            defName,
		DailyOperationalStaff: defName,
		MerchantTypeString:    "REGULAR",
		UpgradeStatus:         "PENDING_MANAGE",
		ProductType:           "PHYSIC",
		IsActive:              true,
		Maps: model.Maps{
			ID:           "MAPS202109063345",
			RelationID:   "ADDR2021030506423",
			RelationName: "b2c_merchant",
			Label:        "Label",
			Latitude:     0,
			Longitude:    0,
		},
	}
)

func generateRepoResult(data repo.ResultRepository) <-chan repo.ResultRepository {
	output := make(chan repo.ResultRepository)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func generateServiceResult(data serviceModel.ServiceResult) <-chan serviceModel.ServiceResult {
	output := make(chan serviceModel.ServiceResult)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestAddMerchant(t *testing.T) {
	defInput2 := defInput
	defInput2.BusinessType = "anything"

	input3 := defInput
	input3.BankID = int32(1)
	input4 := defInput
	input4.Documents = []model.B2CMerchantDocumentInput{
		{
			DocumentValue: "some value",
		},
	}
	input5 := defInput
	input5.Documents = []model.B2CMerchantDocumentInput{
		{
			DocumentType:  "KTP-File",
			DocumentValue: "some value",
		},
	}

	var testDataAddMerchant = []struct {
		name                       string
		wantError                  bool
		input                      *model.B2CMerchantCreateInput
		userAttr                   *model.MerchantUserAttribute
		memberRepoResult           memberRepo.ResultRepository
		repoResult                 repo.ResultRepository
		merchantAddressRepoResult  repo.ResultRepository
		merchantBankRepoResult     repo.ResultRepository
		merchantDocumentRepoResult repo.ResultRepository
		checkExists                repo.ResultRepository
		insertDocument             repo.ResultRepository
		updateDocument             repo.ResultRepository
		checkDocParam              repo.ResultRepository
	}{
		{
			name:             "Test Add Merchant #1", // all passed
			input:            &defInput,
			userAttr:         defUserAttr,
			memberRepoResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoResult:       repo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defaultMerchantID}},
		},
		{
			name:             "Test Add Merchant #2",
			input:            &defInput2,
			userAttr:         defUserAttr,
			memberRepoResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoResult:       repo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defaultMerchantID}},
			wantError:        true,
		},
		{
			name:             "Test Add Merchant #6", // error load member
			input:            &defInput,
			userAttr:         defUserAttr,
			memberRepoResult: memberRepo.ResultRepository{},
			wantError:        true,
		},
		{
			name:             "Test Add Merchant #7", // error merchant exists
			input:            &defInput,
			userAttr:         defUserAttr,
			memberRepoResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			checkExists:      repo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defaultMerchantID}},
			wantError:        true,
		},
		{
			name:             "Test Add Merchant #8", // error processing addUpdate
			input:            &defInput,
			userAttr:         defUserAttr,
			memberRepoResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoResult:       repo.ResultRepository{Error: errDefault},
			wantError:        true,
		},
		{
			name:             "Test Add Merchant #9", // error validating document
			input:            &input4,
			userAttr:         defUserAttr,
			memberRepoResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			wantError:        true,
		},
	}

	for _, tc := range testDataAddMerchant {
		memberRepoMock := mockMemberRepo.MemberRepository{}
		repoMock := mockMerchantRepo.MerchantRepository{}
		merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
		merchantBankRepoMock := mockMerchantRepo.MerchantBankRepository{}
		merchantDocRepoMock := mockMerchantRepo.MerchantDocumentRepository{}
		mockDB, sqlMock, _ := sqlMock.New()
		sqlMock.ExpectBegin()

		defer mockDB.Close()

		svcRepo := localConfig.ServiceRepository{
			MemberRepository:          &memberRepoMock,
			MerchantRepository:        &repoMock,
			MerchantAddressRepository: &merchantAddressRepoMock,
			Repository:                &repository.Repository{WriteDB: mockDB},
		}
		publisher := serviceMock.QPublisher{}
		merchantService := serviceMock.MerchantServices{}
		svcShared := localConfig.ServiceShared{
			QPublisher:      &publisher,
			MerchantService: &merchantService,
		}
		localQuery := localConfig.ServiceQuery{}
		tokenGen := mockToken.AccessTokenGenerator{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		ctxReq := context.Background()

		merchantBankRepoMock.On("FindActiveMerchantBankByID", mock.Anything, mock.Anything).Return(generateRepoResult(tc.merchantBankRepoResult))
		memberRepoMock.On("Load", mock.Anything, mock.Anything).Return(sharedMock.MemberRepoResult(tc.memberRepoResult))
		merchantDocRepoMock.On("FindMerchantDocumentByParam", mock.Anything, mock.Anything).Return(generateRepoResult(tc.checkDocParam))

		merchantDocRepoMock.On("InsertNewMerchantDocument", mock.Anything, mock.Anything).Return(generateRepoResult(tc.insertDocument))
		merchantDocRepoMock.On("UpdateMerchantDocument", mock.Anything, mock.Anything).Return(generateRepoResult(tc.updateDocument))

		repoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.checkExists)
		repoMock.On("FindMerchantByEmail", mock.Anything, mock.Anything).Return(tc.checkExists)
		repoMock.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(tc.checkExists)
		repoMock.On("FindMerchantByName", mock.Anything, mock.Anything).Return(tc.checkExists)
		repoMock.On("FindMerchantBySlug", mock.Anything, mock.Anything).Return(tc.checkExists)
		repoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
		merchantAddressRepoMock.On("AddUpdateAddressMaps", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))

		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
		merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		ucResult := <-m.AddMerchant(ctxReq, tc.input, tc.userAttr)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
			sqlMock.ExpectRollback()
		} else {
			assert.NoError(t, ucResult.Error)
			sqlMock.ExpectCommit()
		}
	}
}

type testUpgradeData struct {
	name                 string
	wantError            bool
	input                *model.B2CMerchantCreateInput
	userAttr             *model.MerchantUserAttribute
	memberRepoLoadResult memberRepo.ResultRepository
	repoResult           repo.ResultRepository
	repoLoadResult       repo.ResultRepository
	serviceResult        serviceModel.ServiceResult
	sendEmailError       error
}

func TestUpgradeMerchant(t *testing.T) {
	os.Setenv("EMAIL_MERCHANT_UPGRADE_TEMPLATE_ID", "11")
	os.Setenv("EMAIL_BCC_HUNTER", "someemail@bhinneka.com")
	os.Setenv("EMAIL_ATTACHMENT_SLA_MERCHANT_MANAGED", "somefile.pdf")
	os.Setenv("EMAIL_ATTACHMENT_SLA_MERCHANT_ASSOCIATE", "somefile.pdf")
	defInput2 := defInput
	defInput2.UpgradeStatus = "PENDING_ASSOCIATE"

	defInput3 := defInput
	defInput3.MerchantTypeString = "ANYTHING"
	defInput4 := defInput
	defInput4.Documents = []model.B2CMerchantDocumentInput{
		{
			DocumentValue: "",
		},
	}

	var testDataUpgradeMerchant = []testUpgradeData{
		{
			name:                 "Test Upgrade Merchant #1",
			input:                &defInput,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:                 "Test Upgrade Merchant #2",
			input:                &defInput2,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:      "Test Upgrade Merchant #3", // FAILED ON validation
			input:     &defInput3,
			userAttr:  defUserAttr,
			wantError: true,
		},
		{
			name:                 "Test Upgrade Merchant #4", // no row on DB while loading memberRepo
			input:                &defInput,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Error: sql.ErrNoRows},
			wantError:            true,
		},
		{
			name:                 "Test Upgrade Merchant #5", // no row on merchant table
			input:                &defInput,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult:       repo.ResultRepository{Error: sql.ErrNoRows},
			wantError:            true,
		},
		{
			name:                 "Test Upgrade Merchant #6", // merchant not active
			input:                &defInput2,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID: defaultMerchantID,
			}},
			wantError: true,
		},
		{
			name:                 "Test Upgrade Merchant #7", // error processing upgrade
			input:                &defInput,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
			repoResult:    repo.ResultRepository{Error: errDefault},
			wantError:     true,
		},
		{
			name:                 "Test Upgrade Merchant #8",
			input:                &defInput4,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
			repoResult:    repo.ResultRepository{Error: errDefault},
			wantError:     true,
		},
		{
			name:                 "Test Upgrade Merchant #9",
			input:                &defInput,
			userAttr:             defUserAttr,
			memberRepoLoadResult: memberRepo.ResultRepository{Result: memberModel.Member{Email: defEmail}},
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult:  serviceModel.ServiceResult{Result: defEmailContent},
			sendEmailError: errDefault,
		},
	}

	for _, tc := range testDataUpgradeMerchant {
		t.Run(tc.name, func(t *testing.T) {
			memberRepoMock := mockMemberRepo.MemberRepository{}
			repoMock := mockMerchantRepo.MerchantRepository{}
			merchantDocumentRepoMock := mockMerchantRepo.MerchantDocumentRepository{}
			mockDB, sqlMock, _ := sqlMock.New()
			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit()

			defer mockDB.Close()

			svcRepo := localConfig.ServiceRepository{
				MemberRepository:   &memberRepoMock,
				MerchantRepository: &repoMock,
				Repository:         &repository.Repository{WriteDB: mockDB},
			}
			publisher := serviceMock.QPublisher{}
			notifService := serviceMock.NotificationServices{}
			merchantService := serviceMock.MerchantServices{}
			svcShared := localConfig.ServiceShared{
				QPublisher:          &publisher,
				MerchantService:     &merchantService,
				NotificationService: &notifService,
			}
			localQuery := localConfig.ServiceQuery{}
			tokenGen := mockToken.AccessTokenGenerator{}
			m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
			ctxReq := context.Background()

			memberRepoMock.On("Load", mock.Anything, mock.Anything).Return(sharedMock.MemberRepoResult(tc.memberRepoLoadResult))
			repoMock.On("FindMerchantByID", mock.Anything, mock.Anything, mock.Anything).Return(tc.repoLoadResult)
			repoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			merchantService.On("InsertLogMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			notifService.On("GetTemplateByID", mock.Anything, mock.Anything, mock.Anything).Return(generateServiceResult(tc.serviceResult))
			notifService.On("SendEmail", mock.Anything, mock.Anything).Return("", tc.sendEmailError)
			merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			merchantDocumentRepoMock.On("InsertUpdateDocument", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
			merchantDocumentRepoMock.On("UpdateMerchantDocument", mock.Anything, mock.Anything, mock.Anything).Return("", nil)
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "SendEmailMerchantUpgrade").Return(nil)
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantUpdate").Return(nil)

			ucResult := <-m.UpgradeMerchant(ctxReq, tc.input, tc.userAttr)
			if tc.wantError {
				assert.Error(t, ucResult.Error)
			} else {
				assert.NoError(t, ucResult.Error)
			}
		})
	}
}

func TestSelfUpdateMerchant(t *testing.T) {
	os.Setenv("EMAIL_MERCHANT_UPGRADE_TEMPLATE_ID", "11")
	os.Setenv("EMAIL_BCC_HUNTER", "someemail@bhinneka.com")
	os.Setenv("EMAIL_ATTACHMENT_SLA_MERCHANT_MANAGED", "somefile.pdf")

	defInput2 := defInput
	defInput2.UpgradeStatus = "MLM"
	defInput3 := defInput

	var testDataUpgradeMerchant = []struct {
		name                  string
		wantError             bool
		input                 *model.B2CMerchantCreateInput
		userAttr              *model.MerchantUserAttribute
		repoResult            repo.ResultRepository
		repoLoadResult        repo.ResultRepository
		serviceResult         serviceModel.ServiceResult
		serviceSendbirdResult serviceModel.ServiceResult
		UpdateUserSendbirdV4  serviceModel.ServiceResult
	}{
		{
			name:     "Test Update Merchant #1",
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:     "Test Update Merchant #2", // error on upgradeStatus
			input:    &defInput2,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			wantError: true,
		},
		{
			name:           "Test Update Merchant #3", // merchant not found
			input:          &defInput3,
			userAttr:       defUserAttr,
			repoLoadResult: repo.ResultRepository{Error: sql.ErrNoRows},
			wantError:      true,
		},
		{
			name:     "Test Update Merchant #4", // failed save
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			repoResult: repo.ResultRepository{Error: errDefault},
			wantError:  true,
		},
	}

	for _, tc := range testDataUpgradeMerchant {
		t.Run(tc.name, func(t *testing.T) {
			memberRepoMock := mockMemberRepo.MemberRepository{}
			repoMock := mockMerchantRepo.MerchantRepository{}
			merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
			mockDB, _, _ := sqlMock.New()
			defer mockDB.Close()

			svcRepo := localConfig.ServiceRepository{
				MemberRepository:          &memberRepoMock,
				MerchantRepository:        &repoMock,
				MerchantAddressRepository: &merchantAddressRepoMock,
				Repository:                &repository.Repository{WriteDB: mockDB},
			}
			publisher := serviceMock.QPublisher{}
			notifService := serviceMock.NotificationServices{}
			merchantService := serviceMock.MerchantServices{}
			sendbirdService := serviceMock.SendbirdServices{}
			svcShared := localConfig.ServiceShared{
				QPublisher:          &publisher,
				MerchantService:     &merchantService,
				NotificationService: &notifService,
				SendbirdService:     &sendbirdService,
			}
			localQuery := localConfig.ServiceQuery{}
			tokenGen := mockToken.AccessTokenGenerator{}
			m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
			ctxReq := context.Background()
			repoMock.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(tc.repoLoadResult)
			repoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			merchantAddressRepoMock.On("AddUpdateAddressMaps", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			sendbirdService.On("UpdateUserSendbirdV4", mock.Anything, mock.Anything).Return(serviceModel.ServiceResult(tc.UpdateUserSendbirdV4))
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantUpdate").Return(nil)
			merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			ucResult := <-m.SelfUpdateMerchant(ctxReq, tc.input, tc.userAttr)
			if tc.wantError {
				assert.Error(t, ucResult.Error)
			} else {
				assert.NoError(t, ucResult.Error)
			}
		})
	}
}

func TestChangeMerchantName(t *testing.T) {
	defInput2 := defInput
	defInput2.MerchantName = "MLMe.21,e,.d,.e"
	defInput3 := defInput
	defInput4 := defInput
	defInput4.IsActive = false

	var testDataChangeMerchantName = []struct {
		name                  string
		wantError             bool
		input                 *model.B2CMerchantCreateInput
		userAttr              *model.MerchantUserAttribute
		repoResult            repo.ResultRepository //find by name
		repoLoadResult        repo.ResultRepository //find by user id
		repoAddUpdateResult   repo.ResultRepository
		MerchantDocumentRepo  repo.MerchantDocumentRepository
		serviceResult         serviceModel.ServiceResult
		serviceSendbirdResult serviceModel.ServiceResult
		UpdateUserSendbirdV4  serviceModel.ServiceResult
	}{
		{
			name:     "Test Change Merchant Name #1", // success
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:                       defaultMerchantID,
				IsActive:                 true,
				CountUpdateNameAvailable: 1,
				MerchantName:             "Testing",
			}},
			repoResult:    repo.ResultRepository{Result: nil},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:     "Test Change Merchant Name #2", //Count Update = 0
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:                       defaultMerchantID,
				IsActive:                 true,
				CountUpdateNameAvailable: 0,
				MerchantName:             "Testing",
			}},
			repoResult: repo.ResultRepository{Result: nil},
			wantError:  true,
		},
		{
			name:           "Test Change Merchant Name #3", // merchant not found
			input:          &defInput3,
			userAttr:       defUserAttr,
			repoLoadResult: repo.ResultRepository{Error: sql.ErrNoRows},
			wantError:      true,
		},
		{
			name:     "Test Change Merchant Name #4", // merchantName sudah ada
			input:    &defInput3,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:                       defaultMerchantID,
				IsActive:                 true,
				CountUpdateNameAvailable: 1,
				MerchantName:             "Testing",
			}},
			repoResult: repo.ResultRepository{Result: errDefault},
			wantError:  true,
		},
		{
			name:     "Test Change Merchant Name #5", // merchant belum aktif
			input:    &defInput4,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:                       defaultMerchantID,
				IsActive:                 false,
				CountUpdateNameAvailable: 1,
				MerchantName:             "Testing",
			}},
			repoResult: repo.ResultRepository{Result: nil},
			wantError:  true,
		},
		{
			name:     "Test Change Merchant Name #6", // failed save
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:                       defaultMerchantID,
				IsActive:                 false,
				CountUpdateNameAvailable: 1,
				MerchantName:             "Testing",
			}},
			repoResult:    repo.ResultRepository{Result: nil},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
			wantError:     true,
		},
		{
			name:     "Test Change Merchant Name #7", // merchant belum aktif
			input:    &defInput4,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:                       defaultMerchantID,
				IsActive:                 false,
				CountUpdateNameAvailable: 1,
				MerchantName:             "Testing",
			}},
			repoResult: repo.ResultRepository{Result: nil},
			wantError:  true,
		},
	}
	for _, tc := range testDataChangeMerchantName {
		t.Run(tc.name, func(t *testing.T) {
			memberRepoMock := mockMemberRepo.MemberRepository{}
			repoMock := mockMerchantRepo.MerchantRepository{}
			merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
			mockDB, _, _ := sqlMock.New()
			defer mockDB.Close()

			svcRepo := localConfig.ServiceRepository{
				MemberRepository:          &memberRepoMock,
				MerchantRepository:        &repoMock,
				MerchantAddressRepository: &merchantAddressRepoMock,
				Repository:                &repository.Repository{WriteDB: mockDB},
			}
			publisher := serviceMock.QPublisher{}
			notifService := serviceMock.NotificationServices{}
			merchantService := serviceMock.MerchantServices{}
			sendbirdService := serviceMock.SendbirdServices{}
			svcShared := localConfig.ServiceShared{
				QPublisher:          &publisher,
				MerchantService:     &merchantService,
				NotificationService: &notifService,
				SendbirdService:     &sendbirdService,
			}
			localQuery := localConfig.ServiceQuery{}
			tokenGen := mockToken.AccessTokenGenerator{}
			m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
			ctxReq := context.Background()
			repoMock.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(tc.repoLoadResult)
			repoMock.On("FindMerchantByName", mock.Anything, mock.Anything).Return(tc.repoResult)
			repoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			merchantAddressRepoMock.On("AddUpdateAddressMaps", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			sendbirdService.On("UpdateUserSendbirdV4", mock.Anything, mock.Anything).Return(serviceModel.ServiceResult(tc.UpdateUserSendbirdV4))
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantUpdate").Return(nil)
			merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			ucResult := <-m.ChangeMerchantName(ctxReq, tc.input, tc.userAttr)
			if tc.wantError {
				assert.Error(t, ucResult.Error)
			} else {
				assert.NoError(t, ucResult.Error)
			}
		})
	}
}

func TestMerchantUseCaseImpl_SelfUpdateMerchantPartial(t *testing.T) {
	defInput2 := defInput
	defInput2.UpgradeStatus = "MLM"
	defInput3 := defInput
	defInput3.Documents = append(defInput3.Documents)

	var testDataSelfUpdateMerchantPartial = []struct {
		name                  string
		wantError             bool
		input                 *model.B2CMerchantCreateInput
		userAttr              *model.MerchantUserAttribute
		query                 model.B2CMerchantDocumentQueryInput
		repoResult            repo.ResultRepository
		repoLoadResult        repo.ResultRepository
		repoDocumentResult    repo.ResultRepository
		serviceResult         serviceModel.ServiceResult
		serviceSendbirdResult serviceModel.ServiceResult
		UpdateUserSendbirdV4  serviceModel.ServiceResult
	}{
		{
			name:     "Test Self Update Merchant Partial #1",
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:      "Test Self Update Merchant Partial #2", // error on upgradeStatus
			input:     &defInput2,
			userAttr:  defUserAttr,
			wantError: true,
		},
		{
			name:           "Test Self Update Merchant Partial #3", // merchant not found
			input:          &defInput3,
			userAttr:       defUserAttr,
			repoLoadResult: repo.ResultRepository{Error: sql.ErrNoRows},
			wantError:      true,
		},
		{
			name:     "Test Self Update Merchant Partial #4", // failed save
			input:    &defInput,
			userAttr: defUserAttr,
			repoLoadResult: repo.ResultRepository{Result: model.B2CMerchantDataV2{
				ID:       defaultMerchantID,
				IsActive: true,
			}},
			repoResult: repo.ResultRepository{Error: errDefault},
			wantError:  true,
		},
	}

	for _, tc := range testDataSelfUpdateMerchantPartial {
		t.Run(tc.name, func(t *testing.T) {
			memberRepoMock := mockMemberRepo.MemberRepository{}
			repoMock := mockMerchantRepo.MerchantRepository{}
			merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
			merchantDocumentRepoMock := mockMerchantRepo.MerchantDocumentRepository{}
			mockDB, _, _ := sqlMock.New()
			defer mockDB.Close()

			svcRepo := localConfig.ServiceRepository{
				MemberRepository:           &memberRepoMock,
				MerchantRepository:         &repoMock,
				MerchantAddressRepository:  &merchantAddressRepoMock,
				MerchantDocumentRepository: &merchantDocumentRepoMock,
				Repository:                 &repository.Repository{WriteDB: mockDB},
			}
			publisher := serviceMock.QPublisher{}
			notifService := serviceMock.NotificationServices{}
			merchantService := serviceMock.MerchantServices{}
			sendbirdService := serviceMock.SendbirdServices{}
			svcShared := localConfig.ServiceShared{
				QPublisher:          &publisher,
				MerchantService:     &merchantService,
				NotificationService: &notifService,
				SendbirdService:     &sendbirdService,
			}
			localQuery := localConfig.ServiceQuery{}
			tokenGen := mockToken.AccessTokenGenerator{}
			m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
			ctxReq := context.Background()
			repoMock.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(tc.repoLoadResult)
			repoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			merchantDocumentRepoMock.On("GetListMerchantDocument", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.repoResult))
			merchantAddressRepoMock.On("FindAddressMaps", mock.Anything, mock.Anything, "b2c_merchant").Return(sharedMock.MerchantRepoResult(tc.repoResult))
			merchantAddressRepoMock.On("AddUpdateAddressMaps", mock.Anything, mock.Anything).Return(generateRepoResult(tc.repoResult))
			sendbirdService.On("UpdateUserSendbirdV4", mock.Anything, mock.Anything).Return(serviceModel.ServiceResult(tc.UpdateUserSendbirdV4))
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantUpdate").Return(nil)
			merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			ucResult := <-m.SelfUpdateMerchantPartial(ctxReq, tc.input, tc.userAttr)
			if tc.wantError {
				assert.Error(t, ucResult.Error)
			} else {
				assert.NoError(t, ucResult.Error)
			}
		})
	}
}
