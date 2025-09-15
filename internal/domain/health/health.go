package health

type Health struct {
	Status  string
	Version string
}

type HealthCheckResponse struct {
	Status    string            `json:"status"`
	Version   string            `json:"version"`
	Database  string            `json:"database"`
	Redis     string            `json:"redis"`
	RabbitMQ  string            `json:"rabbitmq"`
	Timestamp int64             `json:"timestamp"`
	Details   map[string]string `json:"details,omitempty"`
}
