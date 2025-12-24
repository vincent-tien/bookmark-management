package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark_service" envconfig:"SERVICE_NAME"`
	InstanceId  string `envconfig:"INSTANCE_ID"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
