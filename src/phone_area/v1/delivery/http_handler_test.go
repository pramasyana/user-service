package delivery

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
	"github.com/Bhinneka/user-service/src/phone_area/v1/usecase"
	mockUsecase "github.com/Bhinneka/user-service/src/phone_area/v1/usecase/mocks"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func generateUsecaseResult(data usecase.ResultUseCase) <-chan usecase.ResultUseCase {
	output := make(chan usecase.ResultUseCase, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}
func generateRSA() rsa.PublicKey {
	rsaKeyStr := []byte(`{
		"N": 23878505709275011001875030232071538515964203967156573494867521802079450388886948008082271369423710496363779453133485305931627774487834457009042769535758720756791378543746831338298172749747638731118189688519844565774045831849163943719631452593223983696593952639165081060095120464076010454872879321860268068082034083790845080655986972520335163373073393728599406785153011223249135674295571456022713211411571775501137922528076129664967232987827383734947081333879110886185193559381425341463958849336483352888778970004362658494636962670122014112846334846940650524736472570779432379822550640198830292444437468914079622765433,
		"E": 65537
	}`)
	var rsaKey rsa.PublicKey
	json.Unmarshal(rsaKeyStr, &rsaKey)
	return rsaKey
}

func generateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, &middleware.BearerClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return generateRSA(), nil
	})
}

var defaultResult = []model.PhoneArea{
	{
		CodeArea: "021",
	},
	{
		CodeArea: "022",
	},
}

func TestHTTPDelivery(t *testing.T) {
	var testData = []struct {
		name             string
		wantError        bool
		responseUsecase  usecase.ResultUseCase
		expectStatusCode int
	}{
		{
			name:             "Test #1",
			expectStatusCode: http.StatusOK,
			responseUsecase:  usecase.ResultUseCase{Result: defaultResult},
		},
		{
			name:             "Test #2",
			wantError:        true,
			expectStatusCode: http.StatusOK,
			responseUsecase:  usecase.ResultUseCase{Error: errors.New("some error")},
		},
		{
			name:             "Test #3",
			expectStatusCode: http.StatusInternalServerError,
			responseUsecase:  usecase.ResultUseCase{Result: "some value"},
		},
	}
	for _, tc := range testData {
		mockUC := new(mockUsecase.PhoneAreaUseCase)
		mockUC.On("GetAllPhoneArea", mock.Anything).Return(generateUsecaseResult(tc.responseUsecase))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "https://accounts.bhinneka.com/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		handler := NewHTTPHandler(mockUC)
		handler.MountPhoneArea(e.Group("/v1/phone-area"))
		err := handler.GetAllPhoneArea(c)
		if tc.wantError {
			assert.Error(t, err)
		}
		assert.Equal(t, tc.expectStatusCode, rec.Code)
	}
}
