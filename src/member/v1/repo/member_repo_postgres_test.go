package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	sharedRepository "github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/stretchr/testify/assert"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

const (
	TextUSRID                 = "USR123"
	testEmail                 = "test@bhinneka.com"
	expectedQueryVersion      = `^SELECT "version".*`
	expectedQueryInsertMember = `^INSERT INTO member.*`
	expectedQueryUpdateMember = `^UPDATE "member".*`
	expectedQueryLoad         = `^SELECT .*`
	versionText               = "version"
	profilePictureText        = "profilePicture"
)

var (
	memberField = []string{"email", "firstName", "lastName",
		"gender", "mobile", "phone", "ext", "birthDate",
		"password", "salt", "token",
		"province", "provinceId", "city", "cityId",
		"district", "districtId", "subDistrict", "subDistrictId", "zipCode", "address",
		"jobTitle", "department", "isAdmin", "isStaff", "status",
		"facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
		"lastLogin", "lastBlocked", "created", "lastModified", versionText, "isSync", "isActive",
	}

	memberVal = []model.Member{
		{
			Email:        testEmail,
			FirstName:    "Julius2",
			LastName:     "Bernhard2",
			Gender:       model.StringToGender("MALE"),
			Mobile:       "08119889789",
			Phone:        "0217774373",
			Ext:          "270",
			BirthDate:    time.Now(),
			Password:     "DDeYByACAZZfpewIhnASD7LKJZTDh1GzSe8Mkfwqi9uuA8TLvUOBBK+lJ5O7r1JhQ+x0Oj3uBS8q7TNQYk7vJMdBg==2",
			Salt:         "15000.Jjlds4AmJUNxghLXqnSiVLIWtVOxMGKS1cpAQrWLpAlolz2eXU4qqJE4qoxOWjlNGVV6t/6Usy2sjz==2",
			Token:        "token_test",
			JobTitle:     "Senior Software Developer2",
			Department:   "Development2",
			Status:       model.FgStatus(1),
			SignUpFrom:   "starfish2",
			Created:      time.Now(),
			LastModified: time.Now(),
			Version:      1,
			Address: model.Address{
				CityID:        "091",
				City:          "city_test",
				DistrictID:    "092",
				District:      "district_test",
				SubDistrictID: "093",
				SubDistrict:   "sub_district_test",
				ProvinceID:    "094",
				Province:      "province_test",
				ZipCode:       "11111",
				Address:       "jl_testing",
			},
			LastPasswordModified: time.Now(),
			SocialMedia: model.SocialMedia{
				FacebookConnect: time.Now(),
				GoogleConnect:   time.Now(),
				AppleConnect:    time.Now(),
			},
			IsStaff: true,
		},
	}
)

func setupRepoPostgres(t *testing.T) (*MemberRepoPostgres, sqlMock.Sqlmock) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("error new mock %v", err)
	}
	repoMember := &MemberRepoPostgres{
		Repository: &sharedRepository.Repository{
			ReadDB:  db,
			WriteDB: db,
		},
	}
	return repoMember, mock
}

func closeRepoPostgres(r *MemberRepoPostgres) {
	if r.ReadDB != nil {
		r.ReadDB.Close()
	}
	if r.WriteDB != nil {
		r.WriteDB.Close()
	}
}

func TestMemberRepoPostgres(t *testing.T) {

	t.Run("Test Member Repo LoadMember", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows(memberField).
			AddRow("julius.bernhard2@gmail.com", "Julius2", "Bernhard2",
				model.StringToGender("MALE"), "08119889789", "0217774373", "270", time.Now(),
				"DDeYByACAZZfpewIhnFrdnZTDh1GzSe8Mkfwqi9uuA8TLvUOBBK+lJ5O7r1JhQ+x0Oj3uBS8q7TNQYk7vJMdBg==2", "15000.Jjlds4AmJUNxghLXqnSiVLIWtVOxMGKS1cpAQrWLpAlolz2eXU4qqJE4qoxOWjlNGVV6t/6r3nTcXAb8nF0RKg==2", "",
				"", "", "", "",
				"", "", "", "", "", "",
				"Senior Software Developer2", "Development2",
				false, false, "INACTIVE", "", "", "", "", "", "starfish2",
				time.Now(), time.Now(), time.Now(), time.Now(), 1, false, false,
			)

		query := `SELECT email, "firstName", "lastName",
			gender, mobile, phone, ext, "birthDate",
			password, salt, token,
			province, "provinceId", city, "cityId",
			district, "districtId", "subDistrict", "subDistrictId", "zipCode", address,
			"jobTitle", department,
			"isAdmin", "isStaff", status, "facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
			"lastLogin", "lastBlocked", "created", "lastModified", version, "isSync", "isActive"
		FROM member WHERE id = ?`

		mock.ExpectPrepare(query).ExpectQuery().WithArgs("USR1827715").WillReturnRows(rows)

		r.LoadMember("USR1827715")
	})

	t.Run("Test Member Repo Load", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)
		mf := memberField
		mf = append(mf, profilePictureText)
		mf = append(mf, "mfaEnabled")
		mf = append(mf, "mfaKey")
		mf = append(mf, "lastPasswordModified")
		mf = append(mf, "facebookConnect")
		mf = append(mf, "googleConnect")
		mf = append(mf, "appleConnect")
		rows := sqlMock.NewRows(mf).
			AddRow("julius.bernhard@gmail.com", "Julius", "Bernhard",
				model.StringToGender("MALE"), "08119889788", "0217774375", "279", time.Now(),
				"DDeYByACAZZfpewIhnFrdnZTDh1GzSe8Mkfwqi9uuA8TLvUOBBK+lJ5O7r1JhQ+x0Oj3uBS8q7TNQYk7vJMdBg==",
				"15000.Jjlds4AmJUNxghLXqnSiVLIWtVOxMGKS1cpAQrWLpAlolz2eXU4qqJE4qoxOWjlNGVV6t/6r3nTcXAb8nF0RKg==", "",
				"", "", "", "",
				"", "", "", "", "", "",
				"Senior Software Developer", "Development", false, false, "ACTIVE",
				"", "", "", "", "", "starfish",
				time.Now(), time.Now(), time.Now(), time.Now(), 1, "",
				false, "", time.Now(), time.Now(), time.Now(), time.Now(), false, false,
			)

		query := `SELECT email, "firstName", "lastName",
			gender, mobile, phone, ext, "birthDate",
			password, salt, token,
			province, "provinceId", city, "cityId",
			district, "districtId", "subDistrict", "subDistrictId", "zipCode", address,
			"jobTitle", department, "isAdmin", "isStaff", status, 
			"facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
			"lastLogin", "lastBlocked", "created", "lastModified", version, "profilePicture", 
			"mfaEnabled", "mfaKey", "lastPasswordModified", "facebookConnect", "googleConnect", "appleConnect", "isSync", "isActive" 
		FROM member WHERE id = ?`

		mock.ExpectPrepare(query).ExpectQuery().WithArgs("USR1827716").WillReturnRows(rows)

		memberResult := <-r.Load(context.Background(), "USR1827716")

		assert.NoError(t, memberResult.Error)
		assert.IsType(t, model.Member{}, memberResult.Result)
	})

	t.Run("Test Member Repo FindMaxID", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows(
			[]string{"email", "firstName", "lastName",
				"gender", "mobile", "phone", "ext", "birthDate",
				"password", "salt", "token",
				"province", "provinceId", "city", "cityId",
				"district", "districtId", "subDistrict", "subDistrictId", "zipCode", "address",
				"jobTitle", "department", "isAdmin", "isStaff", "status",
				"facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
				"lastLogin", "lastBlocked", "created", "lastModified", versionText,
			}).
			AddRow("julius.bernhard@gmail.com", "Julius", "Bernhard",
				model.StringToGender("MALE"), "08119889788", "0217774375", "279", "1988-07-09",
				"DDeYByACAZZfpewIhnFrdnZTDh1GzSe8Mkfwqi9uuA8TLvUOBBK+lJ5O7r1JhQ+x0Oj3uBS8q7TNQYk7vJMdBg==",
				"15000.Jjlds4AmJUNxghLXqnSiVLIWtVOxMGKS1cpAQrWLpAlolz2eXU4qqJE4qoxOWjlNGVV6t/6r3nTcXAb8nF0RKg==", "",
				"", "", "", "",
				"", "", "", "", "", "",
				"Senior Software Developer", "Development", false, false, "ACTIVE",
				"", "", "", "", "", "starfish",
				"", "", time.Now(), time.Now(), 1,
			)

		query := `[SELECT 
		COALESCE(MAX(CAST(SUBSTRING(id, 8, 15) AS INTEGER)), 0) as max 
		FROM "member"
		WHERE id LIKE CONCAT('USR', TO_CHAR(NOW(), 'YYmm'), '%')]`

		mock.ExpectPrepare(query).ExpectQuery().WillReturnRows(rows)

		memberResult := <-r.FindMaxID(context.Background())

		assert.Error(t, memberResult.Error)
	})

	t.Run("NEGATIVE_MEMBER_REPO_FIND_MAX_ID_PREPARE_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		query := `^[SELECT COALESCE.*]`
		mock.ExpectPrepare(query).WillReturnError(errors.New("error member repo find max id prepare query"))
		mock.ExpectRollback()

		memberResult := <-r.FindMaxID(context.Background())
		assert.Error(t, memberResult.Error)
	})

	t.Run("Test Member Repo UpdateProfilePicture", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		sqlMock.NewRows(
			[]string{"id", profilePictureText}).
			AddRow(TextUSRID, "")
		query := `UPDATE "member" SET "profilePicture" = \? WHERE "id" = \? RETURNING "id", "profilePicture"`
		mock.ExpectPrepare(query).ExpectQuery().WithArgs(TextUSRID, "")
		mock.ExpectExec(query).WithArgs(TextUSRID, "")

		memberResult := <-r.UpdateProfilePicture(context.Background(),
			model.ProfilePicture{
				ID:             "USR125",
				ProfilePicture: "",
			})

		assert.Error(t, memberResult.Error)
	})

	t.Run("NEGATIVE_UPDATE_PROFILE_PICTURE_QUERY_ROW", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectPrepare(expectedQueryUpdateMember).ExpectQuery().WillReturnError(errors.New("error member repo update profile picture prepare query"))

		memberResult := <-r.UpdateProfilePicture(context.Background(), model.ProfilePicture{ID: "USR198765", ProfilePicture: ""})
		assert.Error(t, memberResult.Error)
	})

	t.Run("POSITIVE_UPDATE_PROFILE_PICTURE", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		usrID := "USR522765"
		rows := sqlMock.NewRows([]string{"id", profilePictureText}).AddRow(usrID, "")
		mock.ExpectPrepare(expectedQueryUpdateMember).ExpectQuery().WillReturnRows(rows)

		memberResult := <-r.UpdateProfilePicture(context.Background(), model.ProfilePicture{ID: usrID, ProfilePicture: ""})
		assert.NoError(t, memberResult.Error)
	})

	t.Run("NEGATIVE_MEMBER_REPO_LOAD_MEMBER_PREPARE_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		query := expectedQueryLoad
		mock.ExpectPrepare(query).WillReturnError(errors.New("error member repo load member prepare query"))

		memberResult := r.LoadMember("USR68402")
		assert.Error(t, memberResult.Error)
	})

	t.Run("NEGATIVE_MEMBER_REPO_LOAD_MEMBER_QUERY_ROW", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		query := expectedQueryLoad
		mock.ExpectPrepare(query).ExpectQuery().WillReturnError(errors.New("error member repo load member prepare query row"))

		memberResult := r.LoadMember("USR69952")
		assert.Error(t, memberResult.Error)
	})

	t.Run("NEGATIVE_MEMBER_REPO_LOAD_PREPARE_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		query := expectedQueryLoad
		mock.ExpectPrepare(query).WillReturnError(errors.New("error member repo load prepare query"))

		memberResult := <-r.Load(context.Background(), "USRERROR1")
		assert.Error(t, memberResult.Error)
	})

	t.Run("NEGATIVE_MEMBER_REPO_LOAD_QUERY_ROW", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		query := expectedQueryLoad
		mock.ExpectPrepare(query).ExpectQuery().WillReturnError(errors.New("error member repo load prepare query"))

		memberResult := <-r.Load(context.Background(), "USRERROR2")
		assert.Error(t, memberResult.Error)
	})
}

func TestMemberRepoPostgresBulkSave(t *testing.T) {
	expectedFindMaxIDQuery := `^SELECT COALESCE.*`

	t.Run("NEGATIVE_BULK_SAVE_MEMBER_ID", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{"max"}).AddRow("1")

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedFindMaxIDQuery).ExpectQuery().WillReturnRows(rows)
		mock.ExpectRollback()

		res := <-r.BulkSave(context.Background(), []model.Member{{ID: "USR9345"}})
		assert.Error(t, res.Error)
	})

	t.Run("NEGATIVE_BULK_SAVE_BEGIN_TX", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin().WillReturnError(errors.New("error bulk save tx begin"))
		res := <-r.BulkSave(context.Background(), []model.Member{})
		assert.Error(t, res.Error)
	})

	t.Run("POSITIVE_BULK_SAVE", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{"max"}).AddRow("3")

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedFindMaxIDQuery).ExpectQuery().WillReturnRows(rows)
		mock.ExpectCommit()

		res := <-r.BulkSave(context.Background(), []model.Member{})
		assert.NoError(t, res.Error)
	})
}

func TestNewMemberRepoPostgres(t *testing.T) {
	t.Run("POSITIVE_NEW_MEMBER_REPO", func(t *testing.T) {
		repo := NewMemberRepoPostgres(&sharedRepository.Repository{})
		assert.NotNil(t, repo)
	})
}

func TestMemberRepoPostgresSave(t *testing.T) {
	t.Run("NEGATIVE_SAVE_QUERY_VERSION", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).WillReturnError(errors.New("error select version"))
		res := <-r.Save(context.Background(), model.Member{})
		assert.Error(t, res.Error)
	})

	t.Run("NEGATIVE_SAVE_EXEC_QUERY_VERSION", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnError(errors.New("error exec query version"))
		mock.ExpectRollback()

		res := <-r.Save(context.Background(), model.Member{})
		assert.Error(t, res.Error)
	})

	t.Run("NEGATIVE_SAVE_MEMBER_VERSION", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(3, testEmail)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectRollback()

		res := <-r.Save(context.Background(), model.Member{IsStaff: true, Version: 1, Email: testEmail})
		assert.Error(t, res.Error)
	})

	t.Run("NEGATIVE_SAVE_MEMBER_EMAIL", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, "invalid_email@bhinneka.com")

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectRollback()

		res := <-r.Save(context.Background(), model.Member{IsStaff: true, Version: 2, Email: testEmail})
		assert.Error(t, res.Error)
	})

	t.Run("NEGATIVE_SAVE_MEMBER_PREPARE_INSERT_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, memberVal[0].Email)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertMember).WillReturnError(errors.New("error insert member prepare query"))
		mock.ExpectRollback()

		res := <-r.Save(context.Background(), memberVal[0])
		assert.Error(t, res.Error)
	})

	t.Run("NEGATIVE_SAVE_MEMBER_INSERT_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, memberVal[0].Email)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertMember).ExpectExec().WillReturnError(errors.New("error insert member exec query"))
		mock.ExpectRollback()

		res := <-r.Save(context.Background(), memberVal[0])
		assert.Error(t, res.Error)
	})

	t.Run("POSITIVE_SAVE_MEMBER", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, memberVal[0].Email)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertMember).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectCommit()

		res := <-r.Save(context.Background(), memberVal[0])
		assert.NoError(t, res.Error)
	})
}

func TestMemberRepoPostgresExec(t *testing.T) {
	t.Run("POSITIVE_EXEC", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()

		tx, err := r.WriteDB.Begin()
		if err != nil {
			t.Fatal(err)
		}

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, memberVal[0].Email)
		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertMember).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectCommit()

		err = r.exec(context.Background(), tx, memberVal[0])
		assert.NoError(t, err)
	})

	t.Run("NEGATIVE_EXEC_PREPARE_QUERY_VERSION", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()

		tx, err := r.WriteDB.Begin()
		if err != nil {
			t.Fatal(err)
		}

		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnError(errors.New("error exec member prepare version query"))
		mock.ExpectPrepare(expectedQueryInsertMember).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectRollback()

		err = r.exec(context.Background(), tx, memberVal[0])
		assert.Error(t, err)
	})

	t.Run("NEGATIVE_EXEC_QUERY_VERSION", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()

		tx, err := r.WriteDB.Begin()
		if err != nil {
			t.Fatal(err)
		}

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(3, testEmail)

		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectRollback()

		err = r.exec(context.Background(), tx, memberVal[0])
		assert.Error(t, err)
	})

	t.Run("NEGATIVE_EXEC_EMAIL", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()

		tx, err := r.WriteDB.Begin()
		if err != nil {
			t.Fatal(err)
		}

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, "invalid_email2@bhinneka.com")

		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectRollback()

		err = r.exec(context.Background(), tx, memberVal[0])
		assert.Error(t, err)
	})

	t.Run("NEGATIVE_EXEC_PREPARE_INSERT_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()

		tx, err := r.WriteDB.Begin()
		if err != nil {
			t.Fatal(err)
		}

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, memberVal[0].Email)

		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertMember).WillReturnError(errors.New("error exec prepare insert query"))
		mock.ExpectRollback()

		err = r.exec(context.Background(), tx, memberVal[0])
		assert.Error(t, err)
	})

	t.Run("NEGATIVE_EXEC_INSERT_QUERY", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()

		tx, err := r.WriteDB.Begin()
		if err != nil {
			t.Fatal(err)
		}

		rows := sqlMock.NewRows([]string{versionText, "email"}).AddRow(1, memberVal[0].Email)

		mock.ExpectPrepare(expectedQueryVersion).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertMember).ExpectExec().WillReturnError(errors.New("error exect insert member query"))
		mock.ExpectRollback()

		err = r.exec(context.Background(), tx, memberVal[0])
		assert.Error(t, err)
	})
}
