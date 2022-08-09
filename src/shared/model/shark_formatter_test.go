package model

import (
	"testing"
)

func TestRestructCorporateLeads(t *testing.T) {
	type args struct {
		pl B2BLeadsCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BLeadsCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateLeads(tt.args.pl)
		})
	}
}

func TestRestructCorporateContactAddress(t *testing.T) {
	type args struct {
		pl B2BContactAddressCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BContactAddressCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateContactAddress(tt.args.pl)
		})
	}
}

func TestRestructCorporateDocument(t *testing.T) {
	type args struct {
		pl B2BDocumentCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BDocumentCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateDocument(tt.args.pl)
		})
	}
}

func TestRestructCorporatePhone(t *testing.T) {
	type args struct {
		pl B2BPhoneCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BPhoneCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporatePhone(tt.args.pl)
		})
	}
}

func TestRestructCorporateContact(t *testing.T) {
	type args struct {
		pl B2BContactCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BContactCDC{
					BirthDate:       "1995-09-09",
					TransactionType: `[{"microsite": "MICROSITE_BELA", "type": "shopcart"}]`,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateContact(tt.args.pl)
		})
	}
}

func TestRestructCorporateAddress(t *testing.T) {
	type args struct {
		pl B2BAddressCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BAddressCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateAddress(tt.args.pl)
		})
	}
}

func TestRestructCorporateAccount(t *testing.T) {
	type args struct {
		pl B2BAccountCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BAccountCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateAccount(tt.args.pl)
		})
	}
}

func TestRestructCorporateAccountContact(t *testing.T) {
	type args struct {
		pl B2BAccountContactCDC
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Case 1: Success",
			args: args{
				pl: B2BAccountContactCDC{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RestructCorporateAccountContact(tt.args.pl)
		})
	}
}
