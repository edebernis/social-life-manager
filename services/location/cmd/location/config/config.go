package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config is a global object that holds all application level config variables.
var Config appConfig

type appConfig struct {
	Debug bool

	API struct {
		HTTPBindAddr string
	}
	Metrics struct {
		BindAddr string
		Path     string
	}
	SQL struct {
		Host            string
		Port            int
		User            string
		Password        string
		DB              string
		ConnMaxIdleTime time.Duration
		ConnMaxLifeTime time.Duration
		MaxIdleConns    int
		MaxOpenConns    int
	}
	JWT struct {
		Algorithm string
		Secret    string
	}
}

// LoadConfig loads configuration variables in Config global object
func LoadConfig() error {
	v := viper.New()

	// setup default values for unset config variables
	setDefaults(v)

	// load all env variables starting with LOC prefix
	v.SetEnvPrefix("loc")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// load all config variables into Config global obhect
	if err := v.Unmarshal(&Config); err != nil {
		return fmt.Errorf("Failed to unmarshal config. %w", err)
	}

	return nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("debug", false)

	v.SetDefault("api.httpBindAddr", ":8080")

	v.SetDefault("metrics.bindAddr", ":2112")
	v.SetDefault("metrics.path", "/metrics")

	v.SetDefault("sql.host", "localhost")
	v.SetDefault("sql.port", 5432)
	v.SetDefault("sql.user", "postgres")
	v.SetDefault("sql.password", "")
	v.SetDefault("sql.db", "postgres")
	v.SetDefault("sql.connMaxIdleTime", time.Second*60)
	v.SetDefault("sql.connMaxLifeTime", time.Second*60)
	v.SetDefault("sql.maxIdleConns", 10)
	v.SetDefault("sql.maxOpenConns", 50)

	v.SetDefault("jwt.algorithm", "HS256")
	v.SetDefault("jwt.secret", "")
}
