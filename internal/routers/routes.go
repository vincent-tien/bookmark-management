package routers

// Routes holds the endpoint paths for the API.
type Routes struct {
	HealthCheck  string // Health check endpoint path
	LinkShorten  string // Link shorten endpoint path
	LinkRedirect string // Link redirect endpoint path
	UserRegister string // Link Users register endpoint path
	AuthLogin    string // AuthLogin is the authentication login endpoint path
	GetProfile   string // GetProfile is the user profile retrieval endpoint path
}

var Endpoints = Routes{
	HealthCheck:  "/health-check",
	LinkShorten:  "/links/shorten",
	LinkRedirect: "/links/redirect/*code",
	UserRegister: "/users/register",
	AuthLogin:    "/users/login",
	GetProfile:   "/self/info",
}
