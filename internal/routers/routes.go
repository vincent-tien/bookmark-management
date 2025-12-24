package routers

// Routes holds the endpoint paths for the API.
type Routes struct {
	HealthCheck string // Health check endpoint path
}

var Endpoints = Routes{
	HealthCheck: "/health-check",
}
