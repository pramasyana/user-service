package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4/zero"
)

func TestModel(t *testing.T) {
	t.Run("Test Merchant Bank Data", func(t *testing.T) {
		pa := B2CMerchantBankData{}
		pa.ID = 1
		pa.BankCode = "BCD"
		pa.BankName = zero.StringFrom("BCA")
		pa.Status = true
		pa.CreatorID = "2"
		pa.CreatorIP = "::1"
		pa.Created = "2019-02-11T03:16:50Z"
		pa.EditorID = "USR180662282"
		pa.EditorIP = "::1"
		pa.LastModified = "2019-02-11T09:05:00Z"

		assert.Equal(t, 1, pa.ID)
		assert.Equal(t, "BCD", pa.BankCode)
		assert.Equal(t, "BCA", pa.BankName.String)
		assert.Equal(t, true, pa.Status)
		assert.Equal(t, "2", pa.CreatorID)
		assert.Equal(t, "::1", pa.CreatorIP)
		assert.Equal(t, "2019-02-11T03:16:50Z", pa.Created)
		assert.Equal(t, "USR180662282", pa.EditorID)
		assert.Equal(t, "::1", pa.EditorIP)
		assert.Equal(t, "2019-02-11T09:05:00Z", pa.LastModified)
	})

	t.Run("Test Parameters Merchant", func(t *testing.T) {
		pa := ParametersMerchantBank{}
		pa.StrPage = "1"
		pa.Page = 1
		pa.StrLimit = "10"
		pa.Limit = 10
		pa.Offset = 1
		pa.Sort = "asc"
		pa.OrderBy = "id"
		pa.Status = "true"

		assert.Equal(t, "1", pa.StrPage)
		assert.Equal(t, 1, pa.Page)
		assert.Equal(t, "10", pa.StrLimit)
		assert.Equal(t, 10, pa.Limit)
		assert.Equal(t, 1, pa.Offset)
		assert.Equal(t, "asc", pa.Sort)
		assert.Equal(t, "id", pa.OrderBy)
		assert.Equal(t, "true", pa.Status)
	})

}
