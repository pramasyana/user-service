package helper

import (
	"testing"

	goString "github.com/Bhinneka/golib/string"
	"github.com/stretchr/testify/assert"
)

type baseTest struct {
	name      string
	input     string
	wantError bool
}

var testDataBirthdate = []baseTest{
	{
		name:  "date #1",
		input: "02-10-1986",
	},
	{
		name:  "date #2",
		input: "02-10-0000",
	},
}

var testDataAlphanumeric = []baseTest{
	{
		name:  "alphanum #1",
		input: "pk1a kj",
	},
	{
		name:  "alphanum #2",
		input: "pk1a $%$",
	},
}

var testDataEmail = []baseTest{
	{
		name:      "email #1",
		input:     "pian.mutakin@bhinneka.com",
		wantError: false,
	},
	{
		name:      "email #2",
		input:     "pian.mutakin@bhinneka",
		wantError: true,
	},
}

var testDataLatin = []baseTest{
	{
		name:  "latin #1",
		input: "Some time #",
	},
	{
		name:  "latin #2",
		input: "Some {} time",
	},
	{
		name:  "latin #3",
		input: "Some <> time",
	},
}

// to maintain the same result
// this test should be remove when current helper removed
func TestGolibValidation(t *testing.T) {
	for _, tc := range testDataBirthdate {
		assert.Equal(t, IsValidBirthDate(tc.input), goString.IsValidBirthDate(tc.input))
	}

	for _, tc := range testDataAlphanumeric {
		assert.Equal(t, ValidateAlphaNumericInput(tc.input), goString.ValidateAlphaNumericInput(tc.input))
	}

	for _, tc := range testDataEmail {
		err1 := ValidateEmail(tc.input)
		err2 := goString.ValidateEmail(tc.input)

		if tc.wantError {
			assert.Error(t, err1)
			assert.Error(t, err2)
		} else {
			assert.NoError(t, err1)
			assert.NoError(t, err2)
		}
	}

	for _, tc := range testDataLatin {
		assert.Equal(t, ValidateLatinOnlyExcepTag(tc.input), goString.ValidateLatinOnlyExcepTag(tc.input))
		assert.Equal(t, ValidateLatinOnlyExcepTagCurly(tc.input), goString.ValidateLatinOnlyExcepTagCurly(tc.input))
	}
}
