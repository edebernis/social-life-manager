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
		HTTP struct {
			BindAddr          string
			ReadHeaderTimeout time.Duration
			ReadTimeout       time.Duration
			WriteTimeout      time.Duration
		}
		GRPC struct {
			BindAddr          string
			ConnectionTimeout time.Duration
			AuthScheme        string
		}
		JWT struct {
			Algorithm string
			Secret    string
		}
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
		QueryTimeout    time.Duration
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

	v.SetDefault("api.http.BindAddr", ":8080")
	v.SetDefault("api.http.ReadHeaderTimeout", 20*time.Second)
	v.SetDefault("api.http.ReadTimeout", 1*time.Minute)
	v.SetDefault("api.http.WriteTimeout", 2*time.Minute)

	v.SetDefault("api.grpc.BindAddr", ":9090")
	v.SetDefault("api.grpc.ConnectionTimeout", 2*time.Minute)
	v.SetDefault("api.grpc.AuthScheme", "bearer")

	v.SetDefault("api.jwt.algorithm", "HS256")
	v.SetDefault("api.jwt.secret", "")

	v.SetDefault("metrics.bindAddr", ":2112")
	v.SetDefault("metrics.path", "/metrics")

	v.SetDefault("sql.host", "localhost")
	v.SetDefault("sql.port", 5432)
	v.SetDefault("sql.user", "postgres")
	v.SetDefault("sql.password", "")
	v.SetDefault("sql.db", "postgres")
	v.SetDefault("sql.connMaxIdleTime", 60*time.Second)
	v.SetDefault("sql.connMaxLifeTime", 60*time.Second)
	v.SetDefault("sql.maxIdleConns", 10)
	v.SetDefault("sql.maxOpenConns", 50)
	v.SetDefault("sql.queryTimeout", 5*time.Second)
}
