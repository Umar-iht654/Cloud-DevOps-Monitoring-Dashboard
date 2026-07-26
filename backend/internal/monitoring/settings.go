package monitoring

const (
	// MinCheckIntervalSeconds is the shortest allowed interval between checks for one service.
	MinCheckIntervalSeconds = 45

	// DefaultCheckIntervalSeconds is used when a service does not provide an interval.
	DefaultCheckIntervalSeconds = 60

	// DefaultHealthCheckHistoryLimit is the default number of recent checks returned by the API.
	DefaultHealthCheckHistoryLimit = 25

	// MaxHealthCheckHistoryLimit is the largest number of health checks a client can request.
	MaxHealthCheckHistoryLimit = 100
)
