package config

import "github.com/kelseyhightower/envconfig"

// Config holds the application configuration settings.
// Configuration values are loaded from environment variables with defaults.
type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`                 // Port on which the application runs
	ServiceName string `default:"bookmark_service" envconfig:"SERVICE_NAME"` // Name of the service
	InstanceId  string `envconfig:"INSTANCE_ID"`                             // Unique instance identifier
}

// NewConfig creates a new Config instance by loading values from environment variables.
// It uses the envconfig package to process environment variables with the specified prefixes.
// Returns a pointer to Config and an error if processing fails.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
