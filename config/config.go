package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	DBDriverName     string `envconfig:"DB_DRIVER" default:"sqlite3"`
	DBDataSourceName string `envconfig:"DSN" default:"var/gion.db"`
	RedisAddr        string `envconfig:"REDIS_ADDR" default:"127.0.0.1:6379"`
	HTTPHost         string `envconfig:"HTTP_HOST" default:""`
	HTTPPort         string `envconfig:"HTTP_PORT" default:"8080"`
}

func ReadConfig() (*Config, error) {
	var conf Config
	err := envconfig.Process("gion", &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
