package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	memberRegistrationText = "member-registration"
)

var testDatas = []struct {
	name        string
	key         string
	expectedkey string
	expected    int
}{
	{
		"Test #1",
		memberRegistrationText,
		memberRegistrationText,
		0,
	},
	{
		"Test #2",
		"member-update",
		"member-update",
		1,
	},
	{
		"Test #3",
		"member-activation",
		"member-activation",
		2,
	},
	{
		"Test #4",
		"merchant-registration",
		"merchant-registration",
		3,
	},
	{
		"Test #5",
		"address-create",
		"address-create",
		4,
	},
	{
		"Test #6",
		"address-modify",
		"address-modify",
		5,
	},
	{
		"Test #7",
		"address-primary",
		"address-primary",
		6,
	},
	{
		"Test #8",
		"address-delete",
		"address-delete",
		7,
	},
	{
		"Test #9",
		"",
		memberRegistrationText,
		0,
	},
}

func TestMessageKey(t *testing.T) {
	for _, tc := range testDatas {
		mk := MessageKeyFromString(tc.key)
		assert.Equal(t, int(mk), tc.expected)
		assert.Equal(t, mk.String(), tc.expectedkey)
	}
	var mk MessageKey = 6
	assert.Equal(t, mk.String(), AddressPrimary.String())
	var mkdefault MessageKey = 9
	assert.Equal(t, mkdefault.String(), mkdefault.String())
}
