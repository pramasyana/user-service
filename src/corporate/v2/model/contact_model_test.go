package model

import (
	"testing"

	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/stretchr/testify/assert"
)

func TestRestructCorporateContact(t *testing.T) {
	var inputContact = []struct {
		input    sharedModel.ContactPayloadCDC
		expected string
	}{
		{
			input: sharedModel.ContactPayloadCDC{
				Payload: sharedModel.ContactPayloadData{
					After: sharedModel.B2BContactCDC{
						BirthDate: int32(16463),
					},
				},
			},
			expected: "2015-01-28",
		},
	}

	for _, tc := range inputContact {
		res := sharedModel.RestructCorporateContact(tc.input.Payload.After)
		assert.Equal(t, tc.expected, res.BirthDate.Format(helper.FormatDateDB))
	}
}
