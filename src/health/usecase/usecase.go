package usecase

// ResultUseCase data structure
type ResultUseCase struct {
	Result interface{}
	Error  error
}

// HealthUseCase interface abstraction
type HealthUseCase interface {
	Ping() <-chan ResultUseCase
}
