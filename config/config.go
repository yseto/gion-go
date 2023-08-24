package config

import (
	"context"
	"crypto/rand"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBDriverName     string `envconfig:"DB_DRIVER" default:"sqlite3"`
	DBDataSourceName string `envconfig:"DSN" default:"var/gion.db"`
	RedisAddr        string `envconfig:"REDIS_ADDR" default:"127.0.0.1:6379"`
	HTTPHost         string `envconfig:"HTTP_HOST" default:""`
	HTTPPort         string `envconfig:"HTTP_PORT" default:"8080"`
	JwtSignedKey     string `envconfig:"JWT_SIGNED_KEY" default:""`
	JwtSignedKeyBin  []byte
}

func ReadConfig() (*Config, error) {
	var conf Config
	err := envconfig.Process("gion", &conf)
	if err != nil {
		return nil, err
	}

	if conf.JwtSignedKey == "" {
		// need 256 bits.
		data := make([]byte, 256/8)
		if _, err := rand.Read(data); err != nil {
			return nil, err
		}
		conf.JwtSignedKeyBin = data
	} else {
		conf.JwtSignedKeyBin = []byte(conf.JwtSignedKey)
	}

	return &conf, nil
}

type webConfigKey struct{}

type WebAppConfig struct {
	JwtSignedKeyBin []byte
}

// for web app
func NewContext(ctx context.Context, cfg *Config) context.Context {
	c := &WebAppConfig{
		cfg.JwtSignedKeyBin,
	}
	return context.WithValue(ctx, webConfigKey{}, c)
}

// for web app
func FromContext(ctx context.Context) *WebAppConfig {
	return ctx.Value(webConfigKey{}).(*WebAppConfig)
}
