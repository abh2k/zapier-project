package deployments

const (
	StatusSuccess   = "success"
	StatusFailed    = "failed"
	StatusCancelled = "cancelled"
)

var validStatuses = map[string]struct{}{
	StatusSuccess:   {},
	StatusFailed:    {},
	StatusCancelled: {},
}

func IsValidStatus(status string) bool {
	_, ok := validStatuses[status]
	return ok
}

func ValidStatuses() []string {
	return []string{StatusSuccess, StatusFailed, StatusCancelled}
}
