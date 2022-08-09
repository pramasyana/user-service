package helper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	// alphanumericComaPointDashSpaceRegexp regex for validating alphanumeric, coma, point and space
	alphanumericComaPointDashSpaceRegexp = regexp.MustCompile("^ *([a-zA-Z0-9,.-] ?)+ *$")
	// nonEnglishCharacter regex for validating non  English character
	nonEnglishCharacter = regexp.MustCompile("[^\x00-\x7F]+")

	errorBadFormat    = errors.New("this field value is in bad format")
	errorTooLong      = errors.New("this field value is too long")
	errorShouldNumber = errors.New("this field value should be Number")
)

// PaginationParameters log param struct
type PaginationParameters struct {
	Page     int
	StrPage  string
	Limit    int
	StrLimit string
	Offset   int
}

// IsValidBirthDate validator
// BirthDate should be more than one day ago
// Deprecated: please use golib
func IsValidBirthDate(dob string) bool {
	minimumDob, _ := time.ParseDuration("24h")

	birthDateTime, err := time.Parse("02/01/2006", dob)
	if err != nil {
		return false
	}

	return time.Since(birthDateTime) >= minimumDob
}

// ValidateAlphanumericWithComaDashPointSpace function for validating string to alpha numeric
// Deprecated: no usage
func ValidateAlphanumericWithComaDashPointSpace(str string) error {
	if alphanumericComaPointDashSpaceRegexp.MatchString(str) {
		return nil
	}
	return errors.New("alphanumeric, comma, dash, point and one space between word is allowed")
}

// ValidateAlphaNumericInput function for validating alpha numeric
// Deprecated: please use golib
func ValidateAlphaNumericInput(input string) error {
	check := regexp.MustCompile(`^[A-Za-z0-9.,\\;\\_\\&\\/\-\(\\)\\:\\'\\\ ]+$`).MatchString
	if !check(input) {
		return errorBadFormat
	}
	return nil
}

// ValidateAlphaNumericInputAllowEmpty function for validating alpha numeric and allow empty
// Derecated: no usage
func ValidateAlphaNumericInputAllowEmpty(input string) error {
	check := regexp.MustCompile(`^[A-Za-z0-9.,\\;\\_\\&\\/\-\(\\)\\:\\'\\\ ]*$`).MatchString
	if !check(input) {
		return errorBadFormat
	}
	return nil
}

// ValidateGenderInput function for validating gender
// Deprecated: no usage
func ValidateGenderInput(input string) error {
	check := regexp.MustCompile("^[MFS]+$").MatchString
	if len(input) > 1 {
		err := errors.New("this field value should be M, F or S")
		return err
	}
	if !check(input) {
		return errorBadFormat
	}
	return nil
}

// ValidateBooleanInput function for validating boolean input
// Deprecated: no usage
func ValidateBooleanInput(input string) bool {
	b, e := strconv.ParseBool(input)
	if e != nil {
		return false
	}
	return b
}

// ValidateEmptyInput function for validating empty input
// Deprecated: no usage
func ValidateEmptyInput(input string) error {
	if len(input) == 0 {
		err := errors.New("this field cannot be empty")
		return err
	}
	return nil
}

// ValidateNonEnglishCharacter function for validating string to non english
// Depracated: no usage
func ValidateNonEnglishCharacter(str string) error {
	if nonEnglishCharacter.MatchString(str) {
		return errors.New("character latin only is allowed")
	}

	return nil
}

// ValidatePhoneNumberMaxInput function for validating phone number
func ValidatePhoneNumberMaxInput(input string) error {
	check := regexp.MustCompile("^[0-9 ]*$").MatchString
	if len(input) > 27 {
		return errorTooLong
	}
	if !check(input) {
		return errorShouldNumber
	}
	return nil
}

// ValidateExtPhoneNumberMaxInput function for validating phone number extension
func ValidateExtPhoneNumberMaxInput(input string) error {
	check := regexp.MustCompile("^[0-9 ]*$").MatchString
	if len(input) > 10 {
		return errorTooLong
	}
	if !check(input) {
		return errorShouldNumber
	}
	return nil
}

// ValidateMobileNumberMaxInput function for validating mobile phone number
func ValidateMobileNumberMaxInput(input string) error {
	check := regexp.MustCompile("^[0-9]+$").MatchString
	if len(input) > 13 {
		return errorTooLong
	}
	if !check(input) {
		return errorShouldNumber
	}
	check = regexp.MustCompile(`^[0]+8[1235789]+[0-9]`).MatchString
	if !check(input) {
		err := errors.New("this input mobile number is in bad format")
		return err
	}
	return nil
}

// ValidateNumberOnlyInput function for validating number only
// Deprecated: no usage
func ValidateNumberOnlyInput(input string) error {
	check := regexp.MustCompile("^[0-9]+$").MatchString
	if !check(input) {
		return errorShouldNumber
	}
	return nil
}

// ValidateNumberOnlyInputAllowEmpty function for validating number only and allow empty
// Deprecated: no usage
func ValidateNumberOnlyInputAllowEmpty(input string) error {
	check := regexp.MustCompile("^[0-9]*$").MatchString
	if !check(input) {
		return errorShouldNumber
	}
	return nil
}

// ValidateAlphabeticalOnlyInput function for validating alphabet only
// Deprecated: no usage
func ValidateAlphabeticalOnlyInput(input string) error {
	check := regexp.MustCompile("^[A-Za-z ]+$").MatchString
	if !check(input) {
		err := errors.New("this field value should be alphabetical")
		return err
	}
	return nil
}

// ValidateAlphabeticalOnlyInputAllowEmpty function for validating
// alphabet only and allow empty
// Deprecated: no usage
func ValidateAlphabeticalOnlyInputAllowEmpty(input string) error {
	check := regexp.MustCompile("^[A-Za-z]*$").MatchString
	if !check(input) {
		err := errors.New("this field value should be alphabetical")
		return err
	}
	return nil
}

// ValidateEmail function for validating email
// Derecated: moved to golib
func ValidateEmail(input string) error {
	check := regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,4}$`).MatchString
	if !check(input) {
		err := errors.New("email address is invalid")
		return err
	}
	return nil
}

// IsValidPass func for check valid password
func IsValidPass(pass string) bool {
	if len(pass) == 0 || len(pass) < 6 {
		return false
	}

	var uppercase, lowercase, num, symbol int
	for _, r := range pass {
		if isUppercase(r) {
			uppercase = +1
		} else if isLowercase(r) {
			lowercase = +1
		} else if isNumeric(r) {
			num = +1
		} else if isSpecialCharacter(r) {
			symbol = +1
		}

	}
	return uppercase >= 1 && lowercase >= 1 && num >= 1 && symbol >= 1
}

// ValidatePassword function for validate password
func ValidatePassword(password string) error {
	var err error
	if len(password) < 8 {
		err = errors.New("password cannot less than 8 digit")
		return err
	}

	if !IsValidPass(password) {
		err = errors.New("password contains at least 1 capital letter, 1 lowercase, 1 numeric and 1 special character")
		return err
	}

	return nil
}

// Deprecated: please golib
func isUppercase(r rune) bool {
	return int(r) >= 65 && int(r) <= 90
}

// Deprecated: please use golib
func isLowercase(r rune) bool {
	return int(r) >= 97 && int(r) <= 122
}

// Deprecated: please use golib
func isNumeric(r rune) bool {
	return int(r) >= 48 && int(r) <= 57
}

// ValidateAlphanumeric func for check valid alphanumeric
// Deprecated: no usage
func ValidateAlphanumeric(str string, must bool) bool {
	var uppercase, lowercase, num, symbol int
	for _, r := range str {
		if isUppercase(r) {
			uppercase = +1
		} else if isLowercase(r) {
			lowercase = +1
		} else if isNumeric(r) {
			num = +1
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	if must { //must alphanumeric
		return uppercase >= 1 && lowercase >= 1 && num >= 1
	}

	return uppercase >= 1 || lowercase >= 1 || num >= 1
}

// ValidateHTML function for validating HTML
func ValidateHTML(src string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	return re.ReplaceAllString(src, "")
}

// ValidateLatinOnlyExcepTag func for check valid latin only
// Deprecated: moved to golib
func ValidateLatinOnlyExcepTag(str string) bool {
	var uppercase, lowercase, num, allowed, symbol int
	for _, r := range str {
		if isUppercase(r) {
			uppercase = +1
		} else if isLowercase(r) {
			lowercase = +1
		} else if isNumeric(r) {
			num = +1
		} else if r == 10 || isSpecialCharacter(r) && !isLGCharacter(r) {
			allowed = +1
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	return uppercase >= 1 || lowercase >= 1 || num >= 1 || allowed >= 0
}

// check if character is any of
// code ascii for [space, coma, ., !, ", #, $, %, &, ', (, ), *, +, -, /, :, ;, =, ?, @, [, \, ], ^, _, `, {, |, }, ~]
// original r >= 32 && r <= 47 || r >= 58 && r <= 64 && r != 60 && r != 62 || r >= 91 && r <= 96 || r >= 123 && r <= 126
// Deprecated: moved to golib
func isSpecialCharacter(str rune) bool {
	r := int(str)
	return r >= 32 && r <= 47 || r >= 58 && r <= 64 || r >= 91 && r <= 96 || r >= 123 && r <= 126
}

// check if character is less than (<) or greater than (>)
// Deprecated: moved to golib
func isLGCharacter(str rune) bool {
	return int(str) == 60 || int(str) == 62
}

// check if character is 123 = `{`, 125 = `}`
// Deprecated: moved to golib
func isCurlyBracket(str rune) bool {
	return int(str) == 123 || int(str) == 125
}

// ValidateLatinOnlyExcepTagCurly func for check valid latin only
// Deprecated: moved to golib
func ValidateLatinOnlyExcepTagCurly(str string) bool {
	var uppercase, lowercase, num, allowed, symbol int
	for _, r := range str {
		if isUppercase(r) {
			uppercase = +1
		} else if isLowercase(r) {
			lowercase = +1
		} else if isNumeric(r) {
			num = +1
		} else if isSpecialCharacter(r) && !isCurlyBracket(r) && !isLGCharacter(r) {
			allowed = +1
		} else {
			symbol = +1
		}
	}

	if symbol > 0 {
		return false
	}

	return uppercase >= 1 || lowercase >= 1 || num >= 1 || allowed >= 0
}

// ValidatePagination validates pagination parameters
func ValidatePagination(paging PaginationParameters) (PaginationParameters, error) {
	var err error

	if len(paging.StrPage) > 0 {
		paging.Page, err = strconv.Atoi(paging.StrPage)
		if err != nil || paging.Page <= 0 {
			return paging, fmt.Errorf(ErrorParameterInvalid, "page")
		}
	}

	if len(paging.StrLimit) > 0 {
		paging.Limit, err = strconv.Atoi(paging.StrLimit)
		if err != nil || paging.Limit <= 0 {
			return paging, fmt.Errorf(ErrorParameterInvalid, "limit")
		}
	}

	paging.Offset = (paging.Page - 1) * paging.Limit
	return paging, nil
}
