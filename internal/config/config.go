package config

import (
	"fmt"
	"time"

	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/validate"
	"github.com/spf13/viper"
)

const (
	_defaultConfigurationsPath = "configuration.yaml"
)

type (
	Config struct {
		Logger     *Logger     `validate:"required"`
		Shortener  *Shortener  `validate:"required"`
		HTTPServer *HTTPServer `validate:"required"`
	}

	Logger struct {
		Level int `validate:"min=-4,max=8"`
	}

	Shortener struct {
		Host                       string `validate:"required"`
		CountOfSymbolsInURL        int    `validate:"min=1,max=64"`
		CountOfRetriesCreatingLink int    `validate:"min=1,max=10"`
		Secure                     bool
	}

	HTTPServer struct {
		Port              string        `validate:"required,min=1"`
		MaxHeaderBytes    int           `validate:"gt=0"`
		ReadTimeout       time.Duration `validate:"required"`
		WriteTimeout      time.Duration `validate:"required"`
		ReadHeaderTimeout time.Duration `validate:"required"`
		IdleTimeout       time.Duration `validate:"required"`
	}
)

func New() (*Config, error) {
	vp := viper.New()

	vp.SetConfigFile(_defaultConfigurationsPath)
	if err := vp.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading configurations: %w", err)
	}

	var config Config
	if err := vp.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal configurations: %w", err)
	}

	if err := validate.Struct(config); err != nil {
		return nil, errorspkg.NewValidationError("configuration.New", config, err)
	}

	return &config, nil
}
