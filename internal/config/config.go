package config

import (
	"crypto/rand"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DBDriverName     string `env:"DB_DRIVER" envDefault:"sqlite3"`
	DBDataSourceName string `env:"DSN" envDefault:"var/gion.db"`
	HTTPHost         string `env:"HTTP_HOST" envDefault:""`
	HTTPPort         string `env:"HTTP_PORT" envDefault:"8080"`
	JwtSignedKey     string `env:"JWT_SIGNED_KEY" envDefault:""`
	JwtSignedKeyBin  []byte
}

var conf Config

func init() {
	err := env.Parse(&conf)
	if err != nil {
		panic(err)
	}

	if conf.JwtSignedKey == "" {
		// need 256 bits.
		data := make([]byte, 256/8)
		if _, err := rand.Read(data); err != nil {
			panic("generate rand:" + err.Error())
		}
		conf.JwtSignedKeyBin = data
	} else {
		conf.JwtSignedKeyBin = []byte(conf.JwtSignedKey)
	}
}

func ReadConfig() (*Config, error) {
	return &conf, nil
}
