package service

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/jarcoal/httpmock"
)

func TestNewSendbirdService(t *testing.T) {
	tests := []struct {
		name               string
		sendbirdApiToken   error
		sendbirdServiceUrl error
	}{
		{
			name:               "Case 1: Success",
			sendbirdApiToken:   os.Setenv("SENDBIRD_API_TOKEN", ""),
			sendbirdServiceUrl: os.Setenv("SENDBIRD_SERVICE_URL", ""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewSendbirdService()
		})
	}
}

func TestSendbirdService_CheckUserSenbird(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequest{},
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			bhinneka.MockHTTP(http.MethodGet, `/`, 200, nil)
			p.CheckUserSenbird(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_CheckUserSenbirdV4(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequestV4
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequestV4{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.CheckUserSenbirdV4(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_CreateUserSendbird(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.CreateUserSendbird(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_CreateUserSendbirdV4(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequestV4
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequestV4{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.CreateUserSendbirdV4(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_GetUserSendbird(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.GetUserSendbird(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_GetUserSendbirdV4(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequestV4
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequestV4{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.GetUserSendbirdV4(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_UpdateUserSendbird(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.UpdateUserSendbird(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_UpdateUserSendbirdV4(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequestV4
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequestV4{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.UpdateUserSendbirdV4(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_CreateTokenUserSendbird(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.CreateTokenUserSendbird(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_CreateTokenUserSendbirdV4(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequestV4
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   serviceModel.ServiceResult
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequestV4{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.CreateTokenUserSendbirdV4(tt.args.ctxReq, tt.args.params)
		})
	}
}

func TestSendbirdService_GetTokenUserSendbird(t *testing.T) {
	type fields struct {
		BaseURL  *url.URL
		APIToken string
	}
	type args struct {
		ctxReq context.Context
		params *serviceModel.SendbirdRequest
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				BaseURL:  &url.URL{},
				APIToken: clientIP,
			},
			args: args{
				ctxReq: context.Background(),
				params: &serviceModel.SendbirdRequest{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SendbirdService{
				BaseURL:  tt.fields.BaseURL,
				APIToken: tt.fields.APIToken,
			}
			p.GetTokenUserSendbird(tt.args.ctxReq, tt.args.params)
		})
	}
}
