package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	User     string `default:"ebvn" envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	DbName   string `default:"ebvn_bm" envconfig:"DB_NAME"`
	Port     int    `default:"5432" envconfig:"DB_PORT"`
}

func newConfig(envPrefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(envPrefix, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", c.Host, c.Port, c.User, c.DbName, c.Password)
}
