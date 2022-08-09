package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {

	memberTableTestCases := []*Member{
		{Email: "wuriyanto@bhinneka.com"},
		{Email: "wuriyanto@yahoo.com"},
	}

	t.Run("should return true test IsBhinnekaEmail", func(t *testing.T) {
		m := memberTableTestCases[0]

		assert.True(t, m.IsBhinnekaEmail())
	})

	t.Run("should return false IsBhinnekaEmail", func(t *testing.T) {
		m := memberTableTestCases[1]

		assert.False(t, m.IsBhinnekaEmail())
	})

	t.Run("should return false ValidateGender", func(t *testing.T) {
		_, m := ValidateGender("nonact")
		assert.False(t, m)
	})

	t.Run("should return true ValidateGender Male", func(t *testing.T) {
		_, m := ValidateGender("M")
		assert.True(t, m)
	})

	t.Run("should return true ValidateGender Female", func(t *testing.T) {
		_, m := ValidateGender("F")
		assert.True(t, m)
	})

	t.Run("should return true ValidateGender Secret", func(t *testing.T) {
		_, m := ValidateGender("S")
		assert.True(t, m)
	})

	t.Run("should return false ValidateStatus", func(t *testing.T) {
		_, m := ValidateStatus("nonact")
		assert.False(t, m)
	})

	t.Run("should return true ValidateStatus Active", func(t *testing.T) {
		_, m := ValidateStatus("ACTIVE")
		assert.True(t, m)
	})

	t.Run("should return true ValidateStatus Inactive", func(t *testing.T) {
		_, m := ValidateStatus("INACTIVE")
		assert.True(t, m)
	})

	t.Run("should return true ValidateStatus New", func(t *testing.T) {
		_, m := ValidateStatus("NEW")
		assert.True(t, m)
	})

	t.Run("should return true ValidateStatus Blocked", func(t *testing.T) {
		_, m := ValidateStatus("BLOCKED")
		assert.True(t, m)
	})

}

func TestStatus(t *testing.T) {
	randString := StringToStatus("UNDIFINED")
	getString0 := randString.String()
	assert.EqualValues(t, InactiveString, getString0)

	randString = 5
	getString0 = randString.String()
	assert.EqualValues(t, InactiveString, getString0)

	randString1 := StringToStatus(ActiveString)
	getString1 := randString1.String()
	assert.EqualValues(t, ActiveString, getString1)

	randString2 := StringToStatus(InactiveString)
	getString2 := randString2.String()
	assert.EqualValues(t, InactiveString, getString2)

	randString3 := StringToStatus(BlockedString)
	getString3 := randString3.String()
	assert.EqualValues(t, BlockedString, getString3)

	randString4 := StringToStatus(NewString)
	getString4 := randString4.String()
	assert.EqualValues(t, NewString, getString4)

	genderString1 := StringToGender("UNDIFINED")
	assert.EqualValues(t, 0, genderString1)

	getDolphinString0 := genderString1.GetDolpinGender()
	assert.EqualValues(t, "", getDolphinString0)

	getGenderString0 := genderString1.String()
	assert.EqualValues(t, "", getGenderString0)

	genderString2 := StringToGender(MaleString)
	getGenderString2 := genderString2.String()
	assert.EqualValues(t, MaleString, getGenderString2)

	getDolphinString1 := genderString2.GetDolpinGender()
	assert.EqualValues(t, "M", getDolphinString1)

	genderString3 := StringToGender(FemaleString)
	getGenderString3 := genderString3.String()
	assert.EqualValues(t, FemaleString, getGenderString3)

	getDolphinString3 := genderString3.GetDolpinGender()
	assert.EqualValues(t, "F", getDolphinString3)

	genderString4 := StringToGender(SecretString)
	getGenderString4 := genderString4.String()
	assert.EqualValues(t, SecretString, getGenderString4)

	getDolphinString4 := genderString4.GetDolpinGender()
	assert.EqualValues(t, "S", getDolphinString4)
}
