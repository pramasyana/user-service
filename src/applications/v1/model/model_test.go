package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplicationModel(t *testing.T) {
	applicationsTableTestCases := []*Application{
		{ID: "d817f539-27d3-4789-a695-8491cc0e48f5"},
		{Name: "shark"},
		{Logo: "https://drive.google.com/file/d/1SyT-3quGKXb08S2CnT3qhFIyRf47YUdn/view"},
		{URL: "http://staging.shark.bhinneka.com/"},
	}

	applicationsTableTestNegatif := []*Application{
		{ID: ""},
		{Name: ""},
		{Logo: ""},
		{URL: ""},
	}

	t.Run("Should return true test valid ID", func(t *testing.T) {
		m := applicationsTableTestCases[0]
		assert.True(t, m.ID != "")
	})

	t.Run("should return true test valid Name", func(t *testing.T) {
		m := applicationsTableTestCases[1]
		assert.True(t, m.Name != "")
	})

	t.Run("should return true test valid Logo", func(t *testing.T) {
		m := applicationsTableTestCases[2]
		assert.True(t, m.Logo != "")
	})

	t.Run("should return true test valid URL", func(t *testing.T) {
		m := applicationsTableTestCases[3]
		assert.True(t, m.URL != "")
	})

	t.Run("Should return false test valid ID", func(t *testing.T) {
		m := applicationsTableTestNegatif[0]
		assert.False(t, m.ID != "")
	})

	t.Run("should return false test valid Name", func(t *testing.T) {
		m := applicationsTableTestNegatif[1]
		assert.False(t, m.Name != "")
	})

	t.Run("should return false test valid Logo", func(t *testing.T) {
		m := applicationsTableTestNegatif[2]
		assert.False(t, m.Logo != "")
	})

	t.Run("should return false test valid URL", func(t *testing.T) {
		m := applicationsTableTestNegatif[3]
		assert.False(t, m.URL != "")
	})
}
