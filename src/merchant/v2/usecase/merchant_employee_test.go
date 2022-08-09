package usecase

import (
	"context"
	"errors"
	"testing"

	mocksMerchantUsecase "github.com/Bhinneka/user-service/mocks/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	"github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/stretchr/testify/mock"
	"gopkg.in/guregu/null.v4/zero"
)

func TestMerchantUseCaseImpl_UpdateMerchantEmployee(t *testing.T) {
	type fields struct {
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
	type args struct {
		ctxReq context.Context
		token  string
		params *model.QueryMerchantEmployeeParameters
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{
							Status: zero.NewString("INVITED", true),
						},
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{
					Status: "ACTIVE",
				},
			},
		},
		{
			name: "Case 2: Error token",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{
							Status: zero.NewString("INVITED", true),
						},
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "",
				params: &model.QueryMerchantEmployeeParameters{
					Status: "ACTIVE",
				},
			},
		},
		{
			name: "Case 3: Error FindMerchantByUser",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  errors.New("error"),
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{
							Status: zero.NewString("INVITED", true),
						},
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{
					Status: "ACTIVE",
				},
			},
		},
		{
			name: "Case 4: Error GetMerchantEmployees",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{
							Status: zero.NewString("INVITED", true),
						},
						Error: errors.New("error"),
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{
					Status: "ACTIVE",
				},
			},
		},
		{
			name: "Case 5: Error ChangeStatus",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{
							Status: zero.NewString("INVITED", true),
						},
						Error: nil,
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(errors.New("error"))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{
					Status: "ACTIVE",
				},
			},
		},
		{
			name: "Case 6: Error Vallidate Status",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{
							Status: zero.NewString("INVITED", true),
						},
						Error: nil,
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{
					Status: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantUseCaseImpl{
				Repository:           tt.fields.Repository,
				MerchantRepo:         tt.fields.MerchantRepo,
				MerchantAddressRepo:  tt.fields.MerchantAddressRepo,
				MerchantBankRepo:     tt.fields.MerchantBankRepo,
				MerchantEmployeeRepo: tt.fields.MerchantEmployeeRepo,
				MerchantDocumentRepo: tt.fields.MerchantDocumentRepo,
				MemberRepoRead:       tt.fields.MemberRepoRead,
				UploadService:        tt.fields.UploadService,
				MerchantService:      tt.fields.MerchantService,
				TokenGenerator:       tt.fields.TokenGenerator,
				MemberQueryRead:      tt.fields.MemberQueryRead,
				NotificationService:  tt.fields.NotificationService,
				QueuePublisher:       tt.fields.QueuePublisher,
				SendbirdService:      tt.fields.SendbirdService,
			}
			m.UpdateMerchantEmployee(tt.args.ctxReq, tt.args.token, tt.args.params)
		})
	}
}

func TestMerchantUseCaseImpl_GetMerchantEmployee(t *testing.T) {
	type fields struct {
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
	type args struct {
		ctxReq context.Context
		token  string
		params *model.QueryMerchantEmployeeParameters
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{},
					}))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 2: Error token",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{},
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 3: Error FindMerchantByUser",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  errors.New("error"),
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{},
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 4: Error GetMerchantEmployees",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: model.B2CMerchantEmployeeData{},
						Error:  errors.New("error"),
					}))
					mocksMerchantUsecase.On("ChangeStatus", mock.Anything, mock.Anything).Return(nil)
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantUseCaseImpl{
				Repository:           tt.fields.Repository,
				MerchantRepo:         tt.fields.MerchantRepo,
				MerchantAddressRepo:  tt.fields.MerchantAddressRepo,
				MerchantBankRepo:     tt.fields.MerchantBankRepo,
				MerchantEmployeeRepo: tt.fields.MerchantEmployeeRepo,
				MerchantDocumentRepo: tt.fields.MerchantDocumentRepo,
				MemberRepoRead:       tt.fields.MemberRepoRead,
				UploadService:        tt.fields.UploadService,
				MerchantService:      tt.fields.MerchantService,
				TokenGenerator:       tt.fields.TokenGenerator,
				MemberQueryRead:      tt.fields.MemberQueryRead,
				NotificationService:  tt.fields.NotificationService,
				QueuePublisher:       tt.fields.QueuePublisher,
				SendbirdService:      tt.fields.SendbirdService,
			}
			m.GetMerchantEmployee(tt.args.ctxReq, tt.args.token, tt.args.params)
		})
	}
}

func TestMerchantUseCaseImpl_GetAllMerchantEmployee(t *testing.T) {
	type fields struct {
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
	type args struct {
		ctxReq context.Context
		token  string
		params *model.QueryMerchantEmployeeParameters
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetAllMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: []model.B2CMerchantEmployeeData{},
					}))
					mocksMerchantUsecase.On("GetTotalMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: 1,
					}))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 2: Error token",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetAllMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: []model.B2CMerchantEmployeeData{},
					}))
					mocksMerchantUsecase.On("GetTotalMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: 1,
					}))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 3: Error FindMerchantByUser",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  errors.New("error"),
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetAllMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: []model.B2CMerchantEmployeeData{},
					}))
					mocksMerchantUsecase.On("GetTotalMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: 1,
					}))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 4: Error GetAllMerchantEmployees",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetAllMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: []model.B2CMerchantEmployeeData{},
						Error:  errors.New("error"),
					}))
					mocksMerchantUsecase.On("GetTotalMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: 1,
					}))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
		{
			name: "Case 5: Error GetTotalMerchantEmployees",
			fields: fields{
				MerchantRepo: func() repo.MerchantRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantRepository)
					mocksMerchantUsecase.On("FindMerchantByUser", mock.Anything, mock.Anything).Return(
						repo.ResultRepository{
							Result: model.B2CMerchantDataV2{},
							Error:  nil,
						},
					)
					return mocksMerchantUsecase
				}(),
				MerchantEmployeeRepo: func() repo.MerchantEmployeeRepository {
					mocksMerchantUsecase := new(mocksMerchantUsecase.MerchantEmployeeRepository)
					mocksMerchantUsecase.On("GetAllMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: []model.B2CMerchantEmployeeData{},
					}))
					mocksMerchantUsecase.On("GetTotalMerchantEmployees", mock.Anything, mock.Anything).Return(generateRepoResult(repo.ResultRepository{
						Result: 1,
						Error:  errors.New("error"),
					}))
					return mocksMerchantUsecase
				}(),
			},
			args: args{
				ctxReq: context.Background(),
				token:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOmZhbHNlLCJhdWQiOiJzdHVyZ2VvbiIsImF1dGhvcmlzZWQiOnRydWUsImRpZCI6IlNURzJjYTc2ODljNjE4ZTQ0ZjY5NzFiOTU4NmFiYmM4OTFkIiwiZGxpIjoiV0VCIiwiZW1haWwiOiJzeW5jLnBhc3N3b3JkQHlvcG1haWwuY29tIiwiZXhwIjoxNjQwNzcxNjkyLCJpYXQiOjE2NDA3NjQ0OTIsImlzcyI6ImJoaW5uZWthLmNvbSIsImp0aSI6IjBhNjNhNDc3YjYyYjdlOGIwYTAyN2VmNmMxM2QyZGE3NDQzYmIzMTciLCJtZW1iZXJUeXBlIjoicGVyc29uYWwiLCJzaWduVXBGcm9tIjoic3R1cmdlb24iLCJzdGFmZiI6ZmFsc2UsInN1YiI6IlVTUjIxMTI3MTU5Njg1MjU4In0.dVyDUkaklGj3vH9JPkufmKYQVamNMAQdNM8PcwSYKZOc7pUJj5hP899dPssR6e8AnYsIIsoqd_ZFJlOhKQY0Wvztc1X6L2GYuvmfAUPZMbWsgc7N45CbjOOQ9PEFy0e4Oc0IRqr9egJq-wkyzBSdpc2aMu3sSFR7De1LWdNld6gNmc7nhnKZ23H_eDYirA1y9WkxRBP5F7TNJZSx59rlgAGDgzxJDV-VlhGPsl78TgOygfqfzjHDxxqDlps_RDSAXfi1KNn04j2s-ABdaG2iuIioo14Ia6aErgdlniY0HHvLnG1itLGyTPrnCXH89zXLfeshEXKglV11dL8WZcjfYA",
				params: &model.QueryMerchantEmployeeParameters{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantUseCaseImpl{
				Repository:           tt.fields.Repository,
				MerchantRepo:         tt.fields.MerchantRepo,
				MerchantAddressRepo:  tt.fields.MerchantAddressRepo,
				MerchantBankRepo:     tt.fields.MerchantBankRepo,
				MerchantEmployeeRepo: tt.fields.MerchantEmployeeRepo,
				MerchantDocumentRepo: tt.fields.MerchantDocumentRepo,
				MemberRepoRead:       tt.fields.MemberRepoRead,
				UploadService:        tt.fields.UploadService,
				MerchantService:      tt.fields.MerchantService,
				TokenGenerator:       tt.fields.TokenGenerator,
				MemberQueryRead:      tt.fields.MemberQueryRead,
				NotificationService:  tt.fields.NotificationService,
				QueuePublisher:       tt.fields.QueuePublisher,
				SendbirdService:      tt.fields.SendbirdService,
			}
			m.GetAllMerchantEmployee(tt.args.ctxReq, tt.args.token, tt.args.params)
		})
	}
}
