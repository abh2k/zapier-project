package deployments

import (
	"fmt"
	"math/rand"
	"time"
)

const seedDataLength = 50

func SeedData() []Deployment {
	services := []string{"billing-api", "auth-api", "orders-api", "notifications-api"}
	statuses := ValidStatuses()

	base := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	events := make([]Deployment, 0, seedDataLength)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < seedDataLength; i++ {
		service := services[i%len(services)]
		status := statuses[(i/2)%len(statuses)]
		//
		// Generate a random timestamp between 2025-01-01 and 2025-12-31
		timestamp := base.Add(time.Duration(rng.Intn(365*24)) * time.Hour)
		// Generate a random duration between 60 and 600 seconds
		duration := 60 + rng.Intn(541)
		id := fmt.Sprintf("deploy_%03d", i+1)
		commit := fmt.Sprintf("%06x", rng.Uint32())

		events = append(events, Deployment{
			ID:        id,
			Service:   service,
			Status:    status,
			Duration:  duration,
			Timestamp: timestamp,
			CommitSHA: commit,
		})
	}

	return events
}
