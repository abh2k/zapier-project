package deployments

import "time"

type Deployment struct {
	ID        string    `json:"id"`
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Duration  int       `json:"duration"`
	Timestamp time.Time `json:"timestamp"`
	CommitSHA string    `json:"commit_sha"`
}
