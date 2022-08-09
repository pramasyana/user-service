package delivery

import (
	"errors"
	"math"
	"net/http"

	"github.com/Bhinneka/user-service/src/session/v1/model"
	"github.com/Bhinneka/user-service/src/session/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// HTTPSessionInfoHandler structure
type HTTPSessionInfoHandler struct {
	SessionInfoUseCase usecase.SessionInfoUseCase
}

// NewHTTPHandler function for initialise *HTTPNewSessionInfoHandler
func NewHTTPHandler(sessionInfoUseCase usecase.SessionInfoUseCase) *HTTPSessionInfoHandler {
	return &HTTPSessionInfoHandler{SessionInfoUseCase: sessionInfoUseCase}
}

// MountInfo function for mounting routes
func (h *HTTPSessionInfoHandler) MountInfo(group *echo.Group) {
	group.GET("", h.GetSessionInfoList)
}

// GetSessionInfoList function for get list session data
func (h *HTTPSessionInfoHandler) GetSessionInfoList(c echo.Context) error {
	params := model.ParamList{
		Query:      c.QueryParam("query"),
		StrPage:    c.QueryParam("page"),
		StrLimit:   c.QueryParam("limit"),
		Sort:       c.QueryParam("sort"),
		OrderBy:    c.QueryParam("orderBy"),
		ClientType: c.QueryParam("clientType"),
		Range:      c.QueryParam("range"),
		MemberID:   c.QueryParam("memberId"),
	}

	resultList := <-h.SessionInfoUseCase.GetSessionInfoList(&params)

	if resultList.Error != nil {
		res := shared.NewHTTPResponse(resultList.HTTPStatus, resultList.Error.Error())
		return res.JSON(c)
	}

	response, ok := resultList.Result.(model.SessionInfoList)

	if !ok {
		err := errors.New("result is not list of session info")
		res := shared.NewHTTPResponse(http.StatusBadRequest, err.Error())
		return res.JSON(c)
	}
	totalPage := math.Ceil(float64(response.TotalData) / float64(params.Limit))

	var meta shared.Meta
	meta.Page = params.Page
	meta.Limit = params.Limit
	meta.TotalRecords = response.TotalData
	meta.TotalPages = int(totalPage)

	res := shared.NewHTTPResponse(http.StatusOK, "Success Get All Session Info", response.Data, meta)
	return res.JSON(c)
}
