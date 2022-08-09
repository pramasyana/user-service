package usecase

import (
	"context"
	"reflect"
	"testing"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	mockMemberRepo "github.com/Bhinneka/user-service/mocks/src/member/v1/repo"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	serviceMock "github.com/Bhinneka/user-service/src/service/mocks"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/stretchr/testify/mock"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"

	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/repo"
	mockMerchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo/mocks"

	"github.com/Bhinneka/user-service/src/service"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
)

func TestNewMerchantAddressUseCase(t *testing.T) {
	type args struct {
		repository localConfig.ServiceRepository
		services   localConfig.ServiceShared
	}

	svcRepo := localConfig.ServiceRepository{}
	svcShared := localConfig.ServiceShared{}

	mockNewMerchantAddressUseCase := NewMerchantAddressUseCase(svcRepo, svcShared)
	tests := []struct {
		name string
		args args
		want MerchantAddressUseCase
	}{
		{
			name: "test #1",
			args: args{repository: svcRepo, services: svcShared},
			want: mockNewMerchantAddressUseCase,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMerchantAddressUseCase(tt.args.repository, tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMerchantAddressUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_AddUpdateWarehouseAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq   context.Context
		data     model.WarehouseData
		memberID string
		action   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.AddUpdateWarehouseAddress(tt.args.ctxReq, tt.args.data, tt.args.memberID, tt.args.action); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.AddUpdateWarehouseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FServiceResult(data serviceModel.ServiceResult) <-chan serviceModel.ServiceResult {
	output := make(chan serviceModel.ServiceResult, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestMerchantAddressUseCaseImpl_saveAddressProccess(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	merchantRepoMock := mockMerchantRepo.MerchantRepository{}
	merchantAddressRepoMock := mockMerchantRepo.MerchantAddressRepository{}
	memberRepositoryMock := mockMemberRepo.MemberRepository{}

	mockDB, _, _ := sqlMock.New()
	defer mockDB.Close()

	sRepository := sharedRepo.Repository{WriteDB: mockDB}

	barracudaService := serviceMock.BarracudaServices{}
	activityService := serviceMock.ActivityServices{}
	ctxReq := context.Background()

	type args struct {
		ctxReq        context.Context
		data          model.WarehouseData
		action        string
		memberID      string
		serviceResult serviceModel.ServiceResult
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.WarehouseData
		wantErr bool
	}{
		{
			name: "test #1",
			fields: fields{
				MerchantRepo:        &merchantRepoMock,
				MerchantAddressRepo: &merchantAddressRepoMock,
				MemberRepoRead:      &memberRepositoryMock,
				Repository:          &sRepository,
				BarracudaService:    &barracudaService,
				ActivityService:     &activityService,
			},
			args: args{
				ctxReq:        ctxReq,
				data:          model.WarehouseData{},
				action:        helper.TextUpdate,
				memberID:      "1",
				serviceResult: serviceModel.ServiceResult{},
			},
			want:    model.WarehouseData{},
			wantErr: true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			barracudaService.On("FindZipcode", mock.Anything, mock.Anything).Return(FServiceResult(tt.args.serviceResult))

			got, err := m.saveAddressProccess(tt.args.ctxReq, tt.args.data, tt.args.action, tt.args.memberID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.saveAddressProccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.saveAddressProccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_saveAddressPhoneProccess(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq context.Context
		data   model.WarehouseData
		action string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.WarehouseData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			got, err := m.saveAddressPhoneProccess(tt.args.ctxReq, tt.args.data, tt.args.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.saveAddressPhoneProccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.saveAddressPhoneProccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_validateAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq   context.Context
		data     model.WarehouseData
		memberID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.WarehouseData
		want1   model.WarehouseData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			got, got1, err := m.validateAddress(tt.args.ctxReq, tt.args.data, tt.args.memberID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.validateAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.validateAddress() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("MerchantAddressUseCaseImpl.validateAddress() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_validatePhone(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		data model.PhoneData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.PhoneData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			got, err := m.validatePhone(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.validatePhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.validatePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_validateLocationAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		data model.WarehouseData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.WarehouseData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			got, err := m.validateLocationAddress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.validateLocationAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.validateLocationAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_validateLocationAreaAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		data model.WarehouseData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.WarehouseData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			got, err := m.validateLocationAreaAddress(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.validateLocationAreaAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.validateLocationAreaAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_parseAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq context.Context
		data   model.WarehouseData
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.AddressData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			got, err := m.parseAddress(tt.args.ctxReq, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.parseAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.parseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_UpdatePrimaryWarehouseAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq context.Context
		params model.ParameterPrimaryWarehouse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.UpdatePrimaryWarehouseAddress(tt.args.ctxReq, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.UpdatePrimaryWarehouseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_InsertLogAddressWarehouse(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq  context.Context
		oldData model.WarehouseData
		newData model.WarehouseData
		action  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if err := m.InsertLogAddressWarehouse(tt.args.ctxReq, tt.args.oldData, tt.args.newData, tt.args.action); (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.InsertLogAddressWarehouse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_validateQueryParams(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		params *model.ParameterWarehouse
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if err := m.validateQueryParams(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.validateQueryParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_GetWarehouseAddresses(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq context.Context
		params *model.ParameterWarehouse
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.GetWarehouseAddresses(tt.args.ctxReq, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.GetWarehouseAddresses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_loadMerchant(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq context.Context
		params *model.ParameterWarehouse
		mch    *model.B2CMerchantDataV2
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if err := m.loadMerchant(tt.args.ctxReq, tt.args.params, tt.args.mch); (err != nil) != tt.wantErr {
				t.Errorf("MerchantAddressUseCaseImpl.loadMerchant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_GetDetailWarehouseAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq    context.Context
		addressID string
		memberID  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.GetDetailWarehouseAddress(tt.args.ctxReq, tt.args.addressID, tt.args.memberID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.GetDetailWarehouseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_DeleteWarehouseAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq    context.Context
		addressID string
		memberID  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.DeleteWarehouseAddress(tt.args.ctxReq, tt.args.addressID, tt.args.memberID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.DeleteWarehouseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_FindWarehouseAddress(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq    context.Context
		addressID string
		memberID  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.FindWarehouseAddress(tt.args.ctxReq, tt.args.addressID, tt.args.memberID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.FindWarehouseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerchantAddressUseCaseImpl_GetWarehouseAddressByID(t *testing.T) {
	type fields struct {
		MerchantRepo        repo.MerchantRepository
		MerchantAddressRepo repo.MerchantAddressRepository
		MemberRepoRead      memberRepo.MemberRepository
		Repository          *sharedRepo.Repository
		BarracudaService    service.BarracudaServices
		ActivityService     service.ActivityServices
	}
	type args struct {
		ctxReq     context.Context
		merchantID string
		addressID  string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   <-chan ResultUseCase
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MerchantAddressUseCaseImpl{
				MerchantRepo:        tt.fields.MerchantRepo,
				MerchantAddressRepo: tt.fields.MerchantAddressRepo,
				MemberRepoRead:      tt.fields.MemberRepoRead,
				Repository:          tt.fields.Repository,
				BarracudaService:    tt.fields.BarracudaService,
				ActivityService:     tt.fields.ActivityService,
			}
			if got := m.GetWarehouseAddressByID(tt.args.ctxReq, tt.args.merchantID, tt.args.addressID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MerchantAddressUseCaseImpl.GetWarehouseAddressByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
