package query

// ResultQuery data structure
type ResultQuery struct {
	Result interface{}
	Error  error
}

// HealthQuery interface abstraction
type HealthQuery interface {
	Ping() <-chan ResultQuery
}
