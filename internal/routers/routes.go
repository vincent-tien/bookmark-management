package routers

type Routes struct {
	HealthCheck string
}

var Endpoints = Routes{
	HealthCheck: "/health-check",
}
