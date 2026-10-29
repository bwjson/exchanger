package config

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	env "github.com/caarlos0/env/v11"
)

const path = "config/app.toml"

type Config struct {
	LogLevel slog.Level     `toml:"log_level" env:"LOG_LEVEL"`
	HTTP     HTTPConfig     `toml:"http"`
	Postgres PostgresConfig `toml:"postgres"`
}

type HTTPConfig struct {
	Port int `toml:"port" env:"HTTP_PORT"`
}

type PostgresConfig struct {
	DSN               string        `toml:"-"                   env:"POSTGRES_DSN,required"`
	MaxConns          int32         `toml:"max_conns"           env:"POSTGRES_MAX_CONNS"`
	MinConns          int32         `toml:"min_conns"           env:"POSTGRES_MIN_CONNS"`
	MaxConnLifetime   time.Duration `toml:"max_conn_lifetime"   env:"POSTGRES_MAX_CONN_LIFETIME"`
	MaxConnIdleTime   time.Duration `toml:"max_conn_idle_time"  env:"POSTGRES_MAX_CONN_IDLE_TIME"`
	HealthCheckPeriod time.Duration `toml:"health_check_period" env:"POSTGRES_HEALTH_CHECK_PERIOD"`
	ConnectTimeout    time.Duration `toml:"connect_timeout"     env:"POSTGRES_CONNECT_TIMEOUT"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}
	return cfg, nil
}

func (c Config) ServerAddress() string {
	return net.JoinHostPort("", strconv.Itoa(c.HTTP.Port))
}
