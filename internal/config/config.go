package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Env      string         `yaml:"env"      validate:"required"`
	Server   ServerConfig   `yaml:"server"`
	Log      LogConfig      `yaml:"log"`
	Database DatabaseConfig `yaml:"database"`

	DBPassword string `yaml:"-" env:"DB_PASSWORD"     validate:"required"`
	JWTSecret  string `yaml:"-" env:"AUTH_JWT_SECRET" validate:"min=32"`
}

type ServerConfig struct {
	Name   string       `yaml:"name"   validate:"required"`
	Port   string       `yaml:"port"   validate:"tcp_port"`
	Auth   AuthConfig   `yaml:"auth"`
	Cookie CookieConfig `yaml:"cookie"`
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"  validate:"gt=0"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" validate:"gtfield=AccessTokenTTL"`
}

type CookieConfig struct {
	Secure bool   `yaml:"secure"`
	Domain string `yaml:"domain"`
}

type LogConfig struct {
	Level string `yaml:"level" validate:"oneof=debug info warn error"`
}

type DatabaseConfig struct {
	DBHost    string `yaml:"db_host"     validate:"required"`
	DBPort    int    `yaml:"db_port"     validate:"tcp_port"`
	DBUser    string `yaml:"db_user"     validate:"required"`
	DBName    string `yaml:"db_name"     validate:"required"`
	DBSslMode string `yaml:"db_ssl_mode" validate:"oneof=disable require verify-ca verify-full"`
}

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterAlias("tcp_port", "required,min=1,max=65535")
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		if name, _, _ := strings.Cut(f.Tag.Get("yaml"), ","); name != "" && name != "-" {
			return name
		}
		return f.Tag.Get("env")
	})
	return v
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config %s: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	var cfg Config
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.JWTSecret = os.Getenv("AUTH_JWT_SECRET")

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) validate() error {
	err := validate.Struct(c)
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return err
	}

	errs := make([]error, 0, len(verrs))
	for _, fe := range verrs {
		field := strings.TrimPrefix(fe.Namespace(), "Config.")
		rule := fe.Tag()
		if fe.Param() != "" {
			rule += "=" + fe.Param()
		}
		errs = append(errs, fmt.Errorf("%s: failed rule %q", field, rule))
	}
	return fmt.Errorf("invalid config:\n%w", errors.Join(errs...))
}
