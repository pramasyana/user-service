package query

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
	"github.com/stretchr/testify/assert"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

func TestPhoneAreaQueryPostgresFindAll(t *testing.T) {
	columns := []string{"codeArea", "areaName", "provinceName"}
	query := `SELECT "codeArea", "areaName", "provinceName" FROM phone_area`

	t.Run("Positive Test List Phone Area", func(t *testing.T) {
		db, mock, _ := sqlMock.New()
		defer db.Close()

		rows := sqlMock.NewRows(columns).
			AddRow("0627", "Kota Subulussalam", "Aceh").
			AddRow("0620", "Pangkalan Brandan (Kabupaten Langkat)", "Sumatera Utara").
			AddRow("0753", "Kabupaten Pasaman — Kabupaten Pasaman Barat", "Sumatera Barat")
		mock.ExpectQuery(query).WillReturnRows(rows)

		q := NewPhoneAreaQueryPostgres(db)
		result := <-q.FindAll(context.Background())
		assert.NoError(t, result.Error)

		pa, ok := result.Result.([]model.PhoneArea)
		assert.Equal(t, 3, len(pa))
		assert.True(t, ok)
	})
	t.Run("Negative Test List Phone Area (query error)", func(t *testing.T) {
		db, mock, _ := sqlMock.New()
		defer db.Close()

		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		q := NewPhoneAreaQueryPostgres(db)
		result := <-q.FindAll(context.Background())
		assert.Error(t, result.Error)
	})
	t.Run("Negative Test List Phone Area (scan error)", func(t *testing.T) {
		db, mock, _ := sqlMock.New()
		defer db.Close()

		rows := sqlMock.NewRows(columns).
			AddRow(nil, "Kota Subulussalam", "Aceh")
		mock.ExpectQuery(query).WillReturnRows(rows)

		q := NewPhoneAreaQueryPostgres(db)
		result := <-q.FindAll(context.Background())
		assert.Error(t, result.Error)
	})
}

func TestPhoneAreaQueryPostgresGetTotalPhoneArea(t *testing.T) {
	columns := []string{"count"}
	query := `SELECT .+ FROM phone_area`

	t.Run("Positive Test Total Phone Area", func(t *testing.T) {
		db, mock, _ := sqlMock.New()
		defer db.Close()

		rows := sqlMock.NewRows(columns).AddRow(1)
		mock.ExpectQuery(query).WillReturnRows(rows)

		q := NewPhoneAreaQueryPostgres(db)
		result := <-q.Count(context.Background())
		assert.NoError(t, result.Error)

		pa, ok := result.Result.(model.TotalPhoneArea)
		assert.Equal(t, 1, pa.TotalData)
		assert.True(t, ok)
	})
	t.Run("Negative Test Total Phone Area (query error)", func(t *testing.T) {
		db, mock, _ := sqlMock.New()
		defer db.Close()

		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)

		q := NewPhoneAreaQueryPostgres(db)
		result := <-q.Count(context.Background())
		assert.Error(t, result.Error)
	})
}
