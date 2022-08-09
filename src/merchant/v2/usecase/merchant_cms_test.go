package usecase

import (
	"context"
	"database/sql"
	"os"
	"testing"

	localConfig "github.com/Bhinneka/user-service/config"
	mockToken "github.com/Bhinneka/user-service/src/auth/v1/token/mocks"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	mockMemberQuery "github.com/Bhinneka/user-service/src/member/v1/query/mocks"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	mockMerchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo/mocks"
	serviceMock "github.com/Bhinneka/user-service/src/service/mocks"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	sharedMock "github.com/Bhinneka/user-service/src/shared/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
	"gopkg.in/guregu/null.v4/zero"
)

var testDataReject = []struct {
	name                   string
	wantError              bool
	input                  string
	userAttr               *model.MerchantUserAttribute
	merchantRepoResult     merchantRepo.ResultRepository
	merchantRepoLoadResult merchantRepo.ResultRepository
	serviceResult          serviceModel.ServiceResult
}{
	{
		name:     "Test Reject Merchant Registration #1",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID: defaultMerchantID,
		}},
		serviceResult: serviceModel.ServiceResult{Result: defEmailContent},
	},
	{
		name:                   "Test Reject Merchant Registration #2",
		input:                  defaultMerchantID,
		userAttr:               defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
		serviceResult:          serviceModel.ServiceResult{Result: defEmailContent},
		wantError:              true,
	},
	{
		name:     "Test Reject Merchant Registration #3",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:       defaultMerchantID,
			IsActive: true,
		}},
		serviceResult:      serviceModel.ServiceResult{Result: defEmailContent},
		wantError:          true,
		merchantRepoResult: merchantRepo.ResultRepository{Error: errDefault},
	},
	{
		name:     "Test Reject Merchant Registration #4",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID: defaultMerchantID,
		}},
		serviceResult:      serviceModel.ServiceResult{Result: defEmailContent},
		wantError:          true,
		merchantRepoResult: merchantRepo.ResultRepository{Error: errDefault},
	},
	{
		name:     "Test Reject Merchant Registration #5",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:       defaultMerchantID,
			IsActive: true,
			Status:   model.ActiveString,
		}},
		wantError: true,
	},
}

func TestRejectMerchantRegistration(t *testing.T) {
	for _, tc := range testDataReject {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}

		svcRepo := localConfig.ServiceRepository{
			MerchantRepository: &merchantRepoMock,
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

		merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.merchantRepoLoadResult)
		merchantRepoMock.On("SoftDelete", mock.Anything, mock.Anything).Return(generateRepoResult(tc.merchantRepoResult))
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "SendEmailMerchantRejectRegistration").Return(nil)
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantDelete").Return(nil)
		merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		ucResult := <-m.RejectMerchantRegistration(ctxReq, tc.input, tc.userAttr)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}

func TestDeleteMerchantByID(t *testing.T) {
	del := basicInput{
		name:                "Test Get List Merchant Negative #2",
		wantError:           true,
		input:               defUserID,
		merchantRepoResult:  merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{}},
		merchantRepoResult2: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
	}
	testDataDelete := append(testDataGetListMerchantbyID, del)

	for _, tc := range testDataDelete {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}
		merchantServices := serviceMock.MerchantServices{}
		svcRepo := localConfig.ServiceRepository{
			MerchantRepository: &merchantRepoMock,
		}
		publisher := serviceMock.QPublisher{}
		svcShared := localConfig.ServiceShared{
			MerchantService: &merchantServices,
			QPublisher:      &publisher,
		}
		tokenGen := mockToken.AccessTokenGenerator{}
		ctx := context.Background()

		localQuery := localConfig.ServiceQuery{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.merchantRepoResult)
		merchantRepoMock.On("SoftDelete", mock.Anything, mock.Anything).Return(sharedMock.MerchantRepoResult(tc.merchantRepoResult2))
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantDelete").Return(nil)
		ucResult := <-m.DeleteMerchant(ctx, tc.input, defUserAttr)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}

func TestCreateMerchantFromCMS(t *testing.T) {
	defInput2 := defInput
	defInput2.MerchantName = ""

	var testDataCreateMerchant = []struct {
		name                   string
		wantError              bool
		input                  *model.B2CMerchantCreateInput
		userAttr               *model.MerchantUserAttribute
		memberQueryEmail       memberQuery.ResultQuery
		checkExists            merchantRepo.ResultRepository
		merchantExist          merchantRepo.ResultRepository
		merchantRepoResult     merchantRepo.ResultRepository
		merchantBankRepoResult merchantRepo.ResultRepository
	}{
		{
			name:               "Test Create Merchant #1", // all passed
			input:              &defInput,
			userAttr:           defUserAttr,
			memberQueryEmail:   memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			merchantRepoResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defaultMerchantID}},
		},
		{
			name:             "Test Create Merchant #2", // member not exist
			input:            &defInput,
			userAttr:         defUserAttr,
			memberQueryEmail: memberQuery.ResultQuery{Error: errDefault},
			wantError:        true,
		},
		{
			name:             "Test Create Merchant #3", // bad merchant name
			input:            &defInput2,
			userAttr:         defUserAttr,
			memberQueryEmail: memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			wantError:        true,
		},
		{
			name:             "Test Create Merchant #5", // merchant exist
			input:            &defInput,
			userAttr:         defUserAttr,
			memberQueryEmail: memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			wantError:        true,
			merchantExist:    merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: defaultMerchantID}},
		},
		{
			name:               "Test Create Merchant #6", // failed save
			input:              &defInput,
			userAttr:           defUserAttr,
			memberQueryEmail:   memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			wantError:          true,
			merchantExist:      merchantRepo.ResultRepository{},
			merchantRepoResult: merchantRepo.ResultRepository{Error: errDefault},
		},
		{
			name:             "Test Create Merchant #7", // member not found
			input:            &defInput,
			userAttr:         defUserAttr,
			memberQueryEmail: memberQuery.ResultQuery{Error: sql.ErrNoRows},
			wantError:        true,
		},
		{
			name:             "Test Create Merchant #8", // member not valid
			input:            &defInput,
			userAttr:         defUserAttr,
			memberQueryEmail: memberQuery.ResultQuery{Result: model.ListMerchantBank{}},
			wantError:        true,
		},
	}

	for _, tc := range testDataCreateMerchant {
		memberQueryMock := mockMemberQuery.MemberQuery{}
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}
		mockDB, _, _ := sqlMock.New()
		defer mockDB.Close()

		svcRepo := localConfig.ServiceRepository{
			MerchantRepository: &merchantRepoMock,
			Repository:         &repository.Repository{WriteDB: mockDB},
		}
		publisher := serviceMock.QPublisher{}
		merchantService := serviceMock.MerchantServices{}
		svcShared := localConfig.ServiceShared{
			QPublisher:      &publisher,
			MerchantService: &merchantService,
		}
		localQuery := localConfig.ServiceQuery{
			MemberQueryRead: &memberQueryMock,
		}
		tokenGen := mockToken.AccessTokenGenerator{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		ctxReq := context.Background()
		memberQueryMock.On("FindByEmail", mock.Anything, mock.Anything).Return(sharedMock.MemberQueryResult(tc.memberQueryEmail))
		merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.merchantExist)
		merchantRepoMock.On("FindMerchantByEmail", mock.Anything, mock.Anything).Return(tc.checkExists)
		merchantRepoMock.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(tc.checkExists)
		merchantRepoMock.On("FindMerchantByName", mock.Anything, mock.Anything).Return(tc.checkExists)
		merchantRepoMock.On("FindMerchantBySlug", mock.Anything, mock.Anything).Return(tc.checkExists)

		merchantRepoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.merchantRepoResult))

		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantCreate").Return(nil)
		merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		ucResult := <-m.CreateMerchant(ctxReq, tc.input, tc.userAttr)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}

func TestUpdateMerchantFromCMS(t *testing.T) {
	defInput2 := defInput
	defInput2.UpgradeStatus = "MM"
	os.Setenv("EMAIL_MERCHANT_ACTIVATION", "22")
	os.Setenv("EMAIL_MERCHANT_UPGRADE_APPROVAL", "21")
	os.Setenv("EMAIL_BCC_HUNTER", "email@bhinneka.com")
	defInput3 := defInput
	defInput3.UpgradeStatus = "ACTIVE"
	var rsSendbird serviceModel.SendbirdStringResponse
	rsSendbird.Code = 400201

	var testDataUpdateMerchant = []struct {
		name                     string
		wantError                bool
		input                    *model.B2CMerchantCreateInput
		userAttr                 *model.MerchantUserAttribute
		checkExists              merchantRepo.ResultRepository
		merchantExist            merchantRepo.ResultRepository
		merchantRepoResult       merchantRepo.ResultRepository
		loadMerchantResult       merchantRepo.ResultRepository
		merchantBankRepoResult   merchantRepo.ResultRepository
		serviceResult            serviceModel.ServiceResult
		serviceSendbirdResult    serviceModel.ServiceResult
		createUserSendbirdResult serviceModel.ServiceResult
		errorEmail               error
		memberQueryID            memberQuery.ResultQuery
		responseNotFound         serviceModel.SendbirdErrorResponse
	}{
		{
			name:               "Test Update Merchant #1", // all passed
			input:              &defInput,
			userAttr:           defUserAttr,
			loadMerchantResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{IsActive: false}},
			serviceResult:      serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:      "Test Update Merchant #2", // upgradeStatus not valid
			input:     &defInput2,
			userAttr:  defUserAttr,
			wantError: true,
		},
		{
			name:               "Test Update Merchant #3", // failed update
			input:              &defInput,
			userAttr:           defUserAttr,
			wantError:          true,
			loadMerchantResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{}},
			merchantRepoResult: merchantRepo.ResultRepository{Error: errDefault},
		},
		{
			name:                  "Test Update Merchant #4", // all passed, send email upgrade approval
			input:                 &defInput3,
			userAttr:              defUserAttr,
			loadMerchantResult:    merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: "1", IsActive: true, UpgradeStatus: zero.StringFrom("PENDING_MANAGE")}},
			memberQueryID:         memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			serviceSendbirdResult: serviceModel.ServiceResult{Result: serviceModel.SendbirdStringResponse{}},
			serviceResult:         serviceModel.ServiceResult{Result: defEmailContent},
		},
		{
			name:                  "Test Update Merchant #5", // check userid sendbird
			input:                 &defInput3,
			userAttr:              defUserAttr,
			loadMerchantResult:    merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: "1", IsActive: true, UpgradeStatus: zero.StringFrom("PENDING_MANAGE")}},
			memberQueryID:         memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			serviceSendbirdResult: serviceModel.ServiceResult{Error: errDefault, Result: serviceModel.SendbirdStringResponse{}},
			wantError:             false,
		},
		{
			name:                     "Test Update Merchant #6", // all passed, send email upgrade approval
			input:                    &defInput3,
			userAttr:                 defUserAttr,
			loadMerchantResult:       merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: "1", IsActive: true, UpgradeStatus: zero.StringFrom("PENDING_MANAGE")}},
			memberQueryID:            memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			serviceSendbirdResult:    serviceModel.ServiceResult{Error: errDefault, Result: rsSendbird},
			createUserSendbirdResult: serviceModel.ServiceResult{Result: serviceModel.SendbirdStringResponse{}},
			wantError:                false,
		},

		{
			name:                     "Test Update Merchant #7", // all passed, send email upgrade approval
			input:                    &defInput3,
			userAttr:                 defUserAttr,
			loadMerchantResult:       merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: "1", IsActive: true, UpgradeStatus: zero.StringFrom("PENDING_MANAGE")}},
			memberQueryID:            memberQuery.ResultQuery{Result: memberModel.Member{ID: defUserID}},
			serviceSendbirdResult:    serviceModel.ServiceResult{Error: errDefault, Result: rsSendbird},
			createUserSendbirdResult: serviceModel.ServiceResult{Error: errDefault, Result: serviceModel.SendbirdStringResponse{}},
			wantError:                false,
		},
		{
			name:               "Test Update Merchant #8", // no rows on merchant
			input:              &defInput,
			userAttr:           defUserAttr,
			wantError:          true,
			loadMerchantResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
		},
		{
			name:               "Test Update Merchant #9", // all passed skip send email approval
			input:              &defInput3,
			userAttr:           defUserAttr,
			loadMerchantResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{ID: "1", IsActive: true, UpgradeStatus: zero.StringFrom("PENDING_ASSOCIATE")}},
			serviceResult:      serviceModel.ServiceResult{Result: defEmailContent},
			errorEmail:         errDefault,
		},
	}

	for _, tc := range testDataUpdateMerchant {
		t.Run(tc.name, func(t *testing.T) {
			merchantRepoMock := mockMerchantRepo.MerchantRepository{}
			memberQueryMock := mockMemberQuery.MemberQuery{}
			mockDB, _, _ := sqlMock.New()
			defer mockDB.Close()

			svcRepo := localConfig.ServiceRepository{
				MerchantRepository: &merchantRepoMock,
				Repository:         &repository.Repository{WriteDB: mockDB},
			}
			publisher := serviceMock.QPublisher{}
			merchantService := serviceMock.MerchantServices{}
			notifService := serviceMock.NotificationServices{}
			sendbirdService := serviceMock.SendbirdServices{}

			svcShared := localConfig.ServiceShared{
				QPublisher:          &publisher,
				MerchantService:     &merchantService,
				NotificationService: &notifService,
				SendbirdService:     &sendbirdService,
			}
			localQuery := localConfig.ServiceQuery{
				MemberQueryRead: &memberQueryMock,
			}
			tokenGen := mockToken.AccessTokenGenerator{}

			m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
			ctxReq := context.Background()

			merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.loadMerchantResult)
			merchantRepoMock.On("AddUpdateMerchant", mock.Anything, mock.Anything).Return(generateRepoResult(tc.merchantRepoResult))

			sendbirdService.On("CheckUserSenbird", mock.Anything, mock.Anything).Return(serviceModel.ServiceResult(tc.serviceSendbirdResult))

			merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.loadMerchantResult)

			memberQueryMock.On("FindByID", mock.Anything, mock.Anything).Return(sharedMock.MemberQueryResult(tc.memberQueryID))
			sendbirdService.On("CheckUserSenbird", mock.Anything, mock.Anything).Return(serviceModel.ServiceResult(tc.serviceSendbirdResult))
			sendbirdService.On("CreateUserSendbird", mock.Anything, mock.Anything).Return(serviceModel.ServiceResult(tc.createUserSendbirdResult))

			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "SendEmailActivation").Return(nil)
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "SendEmailApproval").Return(nil)
			publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantUpdate").Return(nil)

			ucResult := <-m.UpdateMerchant(ctxReq, tc.input, tc.userAttr)
			if tc.wantError {
				assert.Error(t, ucResult.Error)
			} else {
				assert.NoError(t, ucResult.Error)
			}
		})

	}
}

var testDataRejectUpgrade = []struct {
	name                   string
	wantError              bool
	input                  string
	userAttr               *model.MerchantUserAttribute
	merchantRepoResult     merchantRepo.ResultRepository
	merchantRepoLoadResult merchantRepo.ResultRepository
	memberQueryResult      memberQuery.ResultQuery
	serviceResult          serviceModel.ServiceResult
	merchantDocRepoResult  merchantRepo.ResultRepository
}{
	{
		name:     "Test Reject Merchant Upgrade #1",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:            defaultMerchantID,
			IsActive:      true,
			UpgradeStatus: zero.StringFrom(model.PendingAssociateString),
		}},
		memberQueryResult:     memberQuery.ResultQuery{Result: memberModel.Member{}},
		merchantDocRepoResult: merchantRepo.ResultRepository{},
	},
	{
		name:                   "Test Reject Merchant Upgrade #2",
		input:                  defaultMerchantID,
		userAttr:               defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Error: sql.ErrNoRows},
		memberQueryResult:      memberQuery.ResultQuery{Result: memberModel.Member{}},
		wantError:              true,
	},
	{
		name:     "Test Reject Merchant Upgrade #3", //  upgrade status is empty
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:       defaultMerchantID,
			IsActive: true,
		}},
		memberQueryResult: memberQuery.ResultQuery{Result: memberModel.Member{}},
		wantError:         true,
	},
	{
		name:     "Test Reject Merchant Upgrade #4", //  upgrade status is active
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:            defaultMerchantID,
			IsActive:      true,
			UpgradeStatus: zero.StringFrom(model.ActiveString),
		}},
		memberQueryResult: memberQuery.ResultQuery{Result: memberModel.Member{}},
		wantError:         true,
	},
	{
		name:     "Test Reject Merchant Upgrade #5",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:            defaultMerchantID,
			IsActive:      true,
			UpgradeStatus: zero.StringFrom(model.PendingAssociateString),
		}},
		wantError:          true,
		memberQueryResult:  memberQuery.ResultQuery{Result: memberModel.Member{}},
		merchantRepoResult: merchantRepo.ResultRepository{Error: errDefault},
	},
	{
		name:     "Test Reject Merchant Upgrade #6",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:            defaultMerchantID,
			IsActive:      true,
			UpgradeStatus: zero.StringFrom(model.PendingAssociateString),
		}},
		wantError:             true,
		memberQueryResult:     memberQuery.ResultQuery{Result: memberModel.Member{}},
		merchantRepoResult:    merchantRepo.ResultRepository{},
		merchantDocRepoResult: merchantRepo.ResultRepository{Error: errDefault},
	},
	{
		name:     "Test Reject Merchant Upgrade #7",
		input:    defaultMerchantID,
		userAttr: defUserAttr,
		merchantRepoLoadResult: merchantRepo.ResultRepository{Result: model.B2CMerchantDataV2{
			ID:            defaultMerchantID,
			IsActive:      true,
			UpgradeStatus: zero.StringFrom(model.PendingManageString),
		}},
		memberQueryResult:     memberQuery.ResultQuery{Result: memberModel.Member{}},
		merchantDocRepoResult: merchantRepo.ResultRepository{},
	},
}

func TestRejectMerchantUpgrade(t *testing.T) {
	for _, tc := range testDataRejectUpgrade {
		merchantRepoMock := mockMerchantRepo.MerchantRepository{}
		merchantDocRepoMock := mockMerchantRepo.MerchantDocumentRepository{}
		memberQueryMock := mockMemberQuery.MemberQuery{}
		mockDB, _, _ := sqlMock.New()
		defer mockDB.Close()
		// sqlMock.ExpectBegin()

		svcRepo := localConfig.ServiceRepository{
			MerchantRepository:         &merchantRepoMock,
			Repository:                 &repository.Repository{WriteDB: mockDB},
			MerchantDocumentRepository: &merchantDocRepoMock,
		}
		publisher := serviceMock.QPublisher{}
		merchantService := serviceMock.MerchantServices{}
		svcShared := localConfig.ServiceShared{
			QPublisher:      &publisher,
			MerchantService: &merchantService,
		}
		localQuery := localConfig.ServiceQuery{
			MemberQueryRead: &memberQueryMock,
		}
		tokenGen := mockToken.AccessTokenGenerator{}
		m := NewMerchantUseCase(svcRepo, svcShared, &tokenGen, localQuery)
		ctxReq := context.Background()

		merchantRepoMock.On("LoadMerchant", mock.Anything, mock.Anything, mock.Anything).Return(tc.merchantRepoLoadResult)
		merchantRepoMock.On("RejectUpgrade", mock.Anything, mock.Anything).Return(generateRepoResult(tc.merchantRepoResult))
		merchantDocRepoMock.On("ResetRejectedDocument", mock.Anything, mock.Anything).Return(generateRepoResult(tc.merchantDocRepoResult))
		memberQueryMock.On("FindByID", mock.Anything, mock.Anything).Return(sharedMock.MemberQueryResult(tc.memberQueryResult))
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "SendEmailMerchantRejectUpgrade").Return(nil)
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "SendEmailAdmin").Return(nil)
		publisher.On("QueueJob", mock.Anything, mock.Anything, mock.Anything, "InsertLogMerchantUpdate").Return(nil)

		merchantService.On("PublishToKafkaUserMerchant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		ucResult := <-m.RejectMerchantUpgrade(ctxReq, tc.input, tc.userAttr, mock.Anything)
		if tc.wantError {
			assert.Error(t, ucResult.Error)
		} else {
			assert.NoError(t, ucResult.Error)
		}
	}
}
