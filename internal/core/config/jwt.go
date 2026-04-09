package core_config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type JWTConfig struct {
	SecretKey       string        `envconfig:"JWT_SECRET_KEY" required:"true"`
	AccessTokenTTL  time.Duration `envconfig:"JWT_ACCESS_TOKEN_TTL" default:"15m"`
	RefreshTokenTTL time.Duration `envconfig:"JWT_REFRESH_TOKEN_TTL" default:"168h"`
}

func NewJWTConfig() (*JWTConfig, error) {
	var cfg JWTConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("process jwt config: %w", err)
	}
	return &cfg, nil
}

func NewJWTConfigMust() *JWTConfig {
	cfg, err := NewJWTConfig()
	if err != nil {
		panic(fmt.Errorf("get jwt config: %w", err))
	}
	return cfg
}
