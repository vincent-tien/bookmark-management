package redis

import "github.com/kelseyhightower/envconfig"

type config struct {
	Addr     string `default:"localhost:6379" envconfig:"REDIS_ADDR"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
