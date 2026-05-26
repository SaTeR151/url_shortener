package app

import (
	"github.com/sater151/url_shortener/internal/config"
	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/validate"
	"github.com/sater151/url_shortener/internal/usecases/shortener"
)

type (
	UsecasesDependencies struct {
		config *config.Config
		repo   *Repo
	}

	Usecases struct {
		shortener *shortener.Manager
	}
)

func NewUsecases(d *UsecasesDependencies) (*Usecases, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewUsecases", d, err)
	}

	shortener, err := shortener.New(&shortener.Dependencies{
		Repo:   d.repo.inMemory,
		Config: d.config.Shortener,
	})
	if err != nil {
		return nil, err
	}

	return &Usecases{
		shortener: shortener,
	}, nil
}
