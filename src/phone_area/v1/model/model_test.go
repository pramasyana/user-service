package model

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	t.Run("Test Phone Area", func(t *testing.T) {
		pa := PhoneArea{}
		pa.CodeArea = "0627"
		pa.AreaName = "Kota Subulussalam"
		pa.ProvinceName = "Aceh"

		assert.Equal(t, "0627", pa.CodeArea)
		assert.Equal(t, "Kota Subulussalam", pa.AreaName)
		assert.Equal(t, "Aceh", pa.ProvinceName)
	})

	t.Run("Test Phone Area Error", func(t *testing.T) {
		pa := PhoneAreaError{}
		pa.ID = "614"
		pa.Message = "Error"

		assert.Equal(t, "614", pa.ID)
		assert.Equal(t, "Error", pa.Message)
	})

	t.Run("Test Phone Area Response", func(t *testing.T) {
		pa1 := PhoneArea{}
		pa1.CodeArea = "0627"
		pa1.AreaName = "Kota Subulussalam 2"
		pa1.ProvinceName = "Aceh"

		pa2 := PhoneArea{}
		pa2.CodeArea = "0627"
		pa2.AreaName = "Kota Subulussalam 2"
		pa2.ProvinceName = "Aceh"

		var phoneArea []PhoneArea

		phoneArea = append(phoneArea, pa1, pa2)

		pa := PhoneAreaResponse{}
		pa.Status = "success"
		pa.Data = phoneArea
		pa.Code = http.StatusOK
		pa.Message = "Get Phone Extension"

		assert.Equal(t, "success", pa.Status)
		assert.Equal(t, http.StatusOK, pa.Code)
		assert.Equal(t, "Get Phone Extension", pa.Message)
	})
}
