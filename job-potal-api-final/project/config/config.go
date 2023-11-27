package config

import (
	"log"

	env "github.com/Netflix/go-env"
)

var cfg Config

type Config struct {
	AppConfig       AppConfig
	DatabaseConfing DatabaseConfing
	KeysPubPri      KeysPubPri
	RadiesConfig    RadiesConfig
}

type AppConfig struct {
	Port string `env:"APP_PORT,required=true"`
}

type DatabaseConfing struct {
	DatabaseConfing1 string `env:"DB_DSN,required=true"`
}

// Host     string `env:"HOST"`
// User     string `env:"USER,required=true"`
// Password string `env:"PASSWORD,required=true"`
// Dbname   string `env:"DBNAME,required=true"`
// Port     string `env:"PORT,required=true"`
// TimeZone string `env:"TIMEZONE,required=true"`
// Sslmode  string `env:"SSLMODE,required=true"`

type KeysPubPri struct {
	Private string `env:"PRIVATE,required=true"`
	Public  string `env:"PUBLIC,required=true"`
}

type RadiesConfig struct {
	Addr     string `env:"ADDR,required=true"`
	Password string `env:"PASSWORD,required=true"`
	DB       int    `env:"DB,required=true"`
}

func init() {
	//var c Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.Fatal(err)
	}

}

func GetConfig() Config {
	return cfg
}
