package routers

// Routes holds the endpoint paths for the API.
type Routes struct {
	HealthCheck  string // Health check endpoint path
	LinkShorten  string // Link shorten endpoint path
	LinkRedirect string // Link redirect endpoint path
}

var Endpoints = Routes{
	HealthCheck:  "/health-check",
	LinkShorten:  "/links/shorten",
	LinkRedirect: "/links/redirect/:code",
}
