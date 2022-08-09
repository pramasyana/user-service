package model

// ResultUseCase data structure
type ResultUseCase struct {
	Result     interface{}
	HTTPStatus int
	Error      error
	Meta       interface{}
}
