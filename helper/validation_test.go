package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidation(t *testing.T) {
	var err error

	validInputString := "Ned Stark"
	invalidInputString := ""

	err = ValidateEmptyInput(validInputString)
	assert.NoError(t, err, err)

	err = ValidateEmptyInput(invalidInputString)
	assert.Error(t, err, err)

	err = ValidateNonEnglishCharacter(validInputString)
	assert.NoError(t, err, err)

	err = ValidateNonEnglishCharacter("汉字")
	assert.Error(t, err, err)

	validNumberInput := "081728828999"
	invalidNumberInput := "981ahbcd"
	emptyNumberInput := ""

	err = ValidateNumberOnlyInput(validNumberInput)
	assert.NoError(t, err, err)

	err = ValidateNumberOnlyInput(invalidNumberInput)
	assert.Error(t, err, err)

	err = ValidateNumberOnlyInputAllowEmpty(validNumberInput)
	assert.NoError(t, err, err)

	err = ValidateNumberOnlyInputAllowEmpty(emptyNumberInput)
	assert.NoError(t, err, err)

	err = ValidateNumberOnlyInputAllowEmpty(invalidNumberInput)
	assert.Error(t, err, err)

	validPhoneNumber := "02129872828"
	validPhoneNumberWithExt := "02129872828 12"
	invalidPhoneNumber := "0218293938111111928272222222"
	invalidPhoneNumberString := "asdjhashdkjashdka"

	err = ValidatePhoneNumberMaxInput(validPhoneNumber)
	assert.NoError(t, err, err)

	err = ValidatePhoneNumberMaxInput(validPhoneNumberWithExt)
	assert.NoError(t, err, err)

	err = ValidatePhoneNumberMaxInput(invalidPhoneNumber)
	assert.Error(t, err, err)

	err = ValidatePhoneNumberMaxInput(invalidPhoneNumberString)
	assert.Error(t, err, err)

	validExtNumber := "0213456789"
	inValidExtNumber := "0213456789123123123"
	inValidExtNumberString := "awertyuio"
	err = ValidateExtPhoneNumberMaxInput(validExtNumber)
	assert.NoError(t, err, err)

	err = ValidateExtPhoneNumberMaxInput(inValidExtNumber)
	assert.Error(t, err, err)

	err = ValidateExtPhoneNumberMaxInput(inValidExtNumberString)
	assert.Error(t, err, err)

	validMobileNumber := "081283187099"
	invalidMobileNumber := "0218293938111111928272"
	invalidMobileNumberString := "adasdhaskd"
	invalidMobileNumberFormat := "0212364789"

	err = ValidateMobileNumberMaxInput(validMobileNumber)
	assert.NoError(t, err, err)

	err = ValidateMobileNumberMaxInput(invalidMobileNumber)
	assert.Error(t, err, err)

	err = ValidateMobileNumberMaxInput(invalidMobileNumberString)
	assert.Error(t, err, err)

	err = ValidateMobileNumberMaxInput(invalidMobileNumberFormat)
	assert.Error(t, err, err)

	validAlphabeticalInput := "AgungDP"
	invalidAlphabeticalInput := "981ahbc"

	err = ValidateAlphabeticalOnlyInput(validAlphabeticalInput)
	assert.NoError(t, err, err)

	err = ValidateAlphabeticalOnlyInput(invalidAlphabeticalInput)
	assert.Error(t, err, err)

	validAlphabeticalEmptyInput := "Prasetyo"

	err = ValidateAlphabeticalOnlyInputAllowEmpty(validAlphabeticalInput)
	assert.NoError(t, err, err)

	err = ValidateAlphabeticalOnlyInputAllowEmpty(validAlphabeticalEmptyInput)
	assert.NoError(t, err, err)

	err = ValidateAlphabeticalOnlyInputAllowEmpty(invalidAlphabeticalInput)
	assert.Error(t, err, err)

	validEmail := "wuriyanto.musobar@gmail.com"
	invalidEmail := "wuriyanto.musobar@gmailcom"

	err = ValidateEmail(validEmail)
	assert.NoError(t, err, err)

	err = ValidateEmail(invalidEmail)
	assert.Error(t, err, err)

	validAlphaNumericDashInput := `nasjd-123`
	inValidAlphaNumericDashInput := `nasjd-汉字`

	err = ValidateAlphanumericWithComaDashPointSpace(validAlphaNumericDashInput)
	assert.NoError(t, err, err)

	err = ValidateAlphanumericWithComaDashPointSpace(inValidAlphaNumericDashInput)
	assert.Error(t, err, err)

	b1 := ValidateBooleanInput("true")
	assert.Equal(t, true, b1)

	b2 := ValidateBooleanInput("false")
	assert.Equal(t, false, b2)

	b3 := ValidateBooleanInput("")
	assert.Equal(t, false, b3)

	validPass := IsValidPass("Blink182!")
	assert.True(t, validPass)

	validPassErrLen := IsValidPass("as!")
	assert.False(t, validPassErrLen)

	validPassErr := IsValidPass("blink182!")
	assert.False(t, validPassErr)

	invalidPassword := RandomPassword("Ab1", 3)
	invalidPassword2 := RandomPassword("ABCD"+"abcd", 8)

	err = ValidatePassword("Som3th!nk")
	assert.Equal(t, nil, err)

	err = ValidatePassword(invalidPassword)
	assert.Error(t, err, err)

	err = ValidatePassword(invalidPassword2)
	assert.Error(t, err, err)

	boolTrue := ValidateAlphanumeric("okesip", false)
	assert.True(t, boolTrue)

	boolTrue2 := ValidateAlphanumeric("Okesip2", true)
	assert.True(t, boolTrue2)

	boolFalse := ValidateAlphanumeric("1FgH^*", false)
	assert.False(t, boolFalse)

	boolFalse = ValidateLatinOnlyExcepTag("<>스칼 k4nAJj1 k0r34")
	assert.False(t, boolFalse)

	boolTrue = ValidateLatinOnlyExcepTag("k4nAJj1 k0r34")
	assert.True(t, boolTrue)

	boolFalse = ValidateLatinOnlyExcepTagCurly("{}<>스칼 k4AJCnj1 k0r34")
	assert.False(t, boolFalse)

	boolTrue = ValidateLatinOnlyExcepTagCurly("k4AJCnj1 k0r34")
	assert.True(t, boolTrue)

	htmlReplace := ValidateHTML("<field>s</field>")
	assert.Equal(t, "s", htmlReplace)
}

func TestAccountValidation(t *testing.T) {

	validInputString := "Ned Stark"
	invalidInputString := ""

	err := ValidateEmptyInput(validInputString)
	assert.NoError(t, err, err)

	err = ValidateEmptyInput(invalidInputString)
	assert.Error(t, err, err)

	validNumberInput := "081728828999"
	invalidNumberInput := "981ahbcd"

	err = ValidateNumberOnlyInput(validNumberInput)
	assert.NoError(t, err, err)

	err = ValidateNumberOnlyInput(invalidNumberInput)
	assert.Error(t, err, err)

	validPhoneNumber := "02129872828"
	validPhoneNumberWithExt := "02129872828 12"
	invalidPhoneNumber := "02182939381111119282722222222"

	err = ValidatePhoneNumberMaxInput(validPhoneNumber)
	assert.NoError(t, err, err)

	err = ValidatePhoneNumberMaxInput(validPhoneNumberWithExt)
	assert.NoError(t, err, err)

	err = ValidatePhoneNumberMaxInput(invalidPhoneNumber)
	assert.Error(t, err, err)

	validMobileNumber := "081283187099"
	invalidMobileNumber := "0218293938111111928272"

	err = ValidateMobileNumberMaxInput(validMobileNumber)
	assert.NoError(t, err, err)

	err = ValidateMobileNumberMaxInput(invalidMobileNumber)
	assert.Error(t, err, err)

	validAlphabeticalInput := "Wuriyanto"
	invalidAlphabeticalInput := "981ahbcde"

	err = ValidateAlphabeticalOnlyInput(validAlphabeticalInput)
	assert.NoError(t, err, err)

	err = ValidateAlphabeticalOnlyInput(invalidAlphabeticalInput)
	assert.Error(t, err, err)

	validGenderInput := "M"
	invalidGenderInput := "mf"
	invalidGenderTypeInput := "g"

	err = ValidateGenderInput(validGenderInput)
	assert.NoError(t, err, err)

	err = ValidateGenderInput(invalidGenderInput)
	assert.Error(t, err, err)

	err = ValidateGenderInput(invalidGenderTypeInput)
	assert.Error(t, err, err)
}

func TestValidateAlphanumeric(t *testing.T) {

	testsEmptyInput := []struct {
		input     string
		wantError bool
	}{
		{
			input:     `dwi'017 / & , bla. - () :\`,
			wantError: false,
		},
		{
			input:     "aahaj 893*",
			wantError: true,
		},
	}

	for _, tt := range testsEmptyInput {

		err := ValidateAlphaNumericInput(tt.input)
		if !tt.wantError {
			assert.NoError(t, err, err)
		} else {
			assert.Error(t, err, err)
		}
	}

	testsAllowEmptyInput := []struct {
		input     string
		wantError bool
	}{
		{
			input:     `prasetyo'002 / & , bla. - () :\`,
			wantError: false,
		},
		{
			input:     "",
			wantError: false,
		},
		{
			input:     "aahaj 893*",
			wantError: true,
		},
	}

	for _, tt := range testsAllowEmptyInput {

		err := ValidateAlphaNumericInputAllowEmpty(tt.input)
		if !tt.wantError {
			assert.NoError(t, err, err)
		} else {
			assert.Error(t, err, err)
		}
	}
}

func TestBirthDateValidator(t *testing.T) {
	validInputs := []string{
		"13/10/1990",
		"13/10/2006",
		"13/01/2007",
		"03/02/2019",
	}

	invalidInputs := []string{
		"13/10/3919",
		"13/04/3919",
		"",
	}

	for _, d := range validInputs {
		valid := IsValidBirthDate(d)
		assert.True(t, valid)
	}

	for _, d := range invalidInputs {
		valid := IsValidBirthDate(d)
		assert.False(t, valid)
	}

}

func TestValidatePagination(t *testing.T) {

	testsEmptyInput := []struct {
		Page, Limit string
		wantError   bool
	}{
		{
			Page:      "1",
			Limit:     "1",
			wantError: false,
		},
		{
			Page:      "sadasd",
			Limit:     "2",
			wantError: true,
		},
		{
			Page:      "1",
			Limit:     "asdsd",
			wantError: true,
		},
	}

	for _, tt := range testsEmptyInput {
		paging := PaginationParameters{
			StrPage:  tt.Page,
			StrLimit: tt.Limit,
		}

		_, err := ValidatePagination(paging)
		if !tt.wantError {
			assert.NoError(t, err, err)
		} else {
			assert.Error(t, err, err)
		}
	}
}
