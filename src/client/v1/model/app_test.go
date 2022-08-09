package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStatus_String(t *testing.T) {
	tests := []struct {
		name string
		s    Status
		want string
	}{
		{
			name: "Case 1: Inactive",
			s:    InActive,
			want: inactive,
		},
		{
			name: "Case 2: Active",
			s:    Active,
			want: active,
		},
		{
			name: "Case 3: Blocked",
			s:    Blocked,
			want: blocked,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringToStatus(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Status
	}{
		{
			name: "Case 1: active",
			args: args{
				s: active,
			},
			want: 1,
		},
		{
			name: "Case 2: inactive",
			args: args{
				s: inactive,
			},
			want: 0,
		},
		{
			name: "Case 3: blocked",
			args: args{
				s: blocked,
			},
			want: 2,
		},
		{
			name: "Case 3: empty",
			args: args{
				s: "",
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringToStatus(tt.args.s)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestClientApp_Authenticate(t *testing.T) {
	type fields struct {
		ID           string
		ClientID     string
		ClientSecret string
		Name         string
		Status       Status
		Created      time.Time
		LastModified time.Time
		Version      int
		Secret       string
	}
	type args struct {
		secret string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "Case 1: Success",
			fields: fields{},
			args: args{
				secret: "secret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientApp{
				ID:           tt.fields.ID,
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				Name:         tt.fields.Name,
				Status:       tt.fields.Status,
				Created:      tt.fields.Created,
				LastModified: tt.fields.LastModified,
				Version:      tt.fields.Version,
				Secret:       tt.fields.Secret,
			}
			c.Authenticate(tt.args.secret)
		})
	}
}

func TestClientApp_IsActive(t *testing.T) {
	type fields struct {
		ID           string
		ClientID     string
		ClientSecret string
		Name         string
		Status       Status
		Created      time.Time
		LastModified time.Time
		Version      int
		Secret       string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Case 1: Success",
			fields: fields{
				Status: Active,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ClientApp{
				ID:           tt.fields.ID,
				ClientID:     tt.fields.ClientID,
				ClientSecret: tt.fields.ClientSecret,
				Name:         tt.fields.Name,
				Status:       tt.fields.Status,
				Created:      tt.fields.Created,
				LastModified: tt.fields.LastModified,
				Version:      tt.fields.Version,
				Secret:       tt.fields.Secret,
			}
			if got := c.IsActive(); got != tt.want {
				t.Errorf("ClientApp.IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
