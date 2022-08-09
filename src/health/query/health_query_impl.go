package query

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Bhinneka/user-service/config/redis"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/health/model"
)

// HealthQueryImpl data structure
type HealthQueryImpl struct {
	name          string
	postgresDB    *sql.DB
	Client        redis.Client
	nsqServiceURL string
}

// NewHealthQueryImpl function for initializing health query implementation
func NewHealthQueryImpl(name string, postgresDB *sql.DB, redisDB redis.Client, nsqBaseURL, nsqTopic string) HealthQuery {
	nsqServiceURL := fmt.Sprintf("http://%s/lookup?topic=%s", nsqBaseURL, nsqTopic)
	return &HealthQueryImpl{name: name, postgresDB: postgresDB, Client: redisDB, nsqServiceURL: nsqServiceURL}
}

// Ping function for checking service
func (q *HealthQueryImpl) Ping() <-chan ResultQuery {

	output := make(chan ResultQuery)

	go func() {
		defer close(output)

		health := &model.Health{}
		//check all dependencies, external services and driver
		q.checkPostgresDB(health)
		q.checkRedis(health)
		q.checkTimeZone(health)

		if health.ErrorCount > 0 {
			health.State = http.StatusServiceUnavailable
		} else {
			health.State = http.StatusOK
		}

		output <- ResultQuery{Result: health}

	}()

	return output
}

//checkPostgresDB function for check Postgres DB State
func (q *HealthQueryImpl) checkTimeZone(health *model.Health) {
	currentTime := time.Now()
	formattedDate := currentTime.Format(time.RFC3339)
	health.Dependencies = append(health.Dependencies, "Timezone : "+formattedDate)

}

//checkPostgresDB function for check Postgres DB State
func (q *HealthQueryImpl) checkPostgresDB(health *model.Health) {
	err := q.postgresDB.Ping()
	if err != nil {
		health.ErrorCount = health.ErrorCount + 1
		health.Dependencies = append(health.Dependencies, fmt.Sprintf("Postgres DB State: %s", err.Error()))
		helper.SendErrorLog(context.Background(), "checkPostgresDB", "check_psql", err, nil)
	} else {
		health.Dependencies = append(health.Dependencies, "Postgres DB State: OK")
	}
}

//checkRedis function for check Redis DB State
func (q *HealthQueryImpl) checkRedis(health *model.Health) {
	pong, err := q.Client.Ping()
	if err != nil {
		health.ErrorCount = health.ErrorCount + 1
		health.Dependencies = append(health.Dependencies, fmt.Sprintf("REDIS State: %s", err.Error()))
		helper.SendErrorLog(context.Background(), "checkRedis", "check_redis", err, nil)
	} else {
		health.Dependencies = append(health.Dependencies, fmt.Sprintf("REDIS State: %s", pong))
	}
}
