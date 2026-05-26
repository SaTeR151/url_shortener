package app

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/sater151/url_shortener/internal/config"
	"github.com/sater151/url_shortener/internal/config/credentials"
	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/logger"
	"github.com/sater151/url_shortener/internal/pkg/validate"
)

type (
	Dependencies struct {
		Config *config.Config `validate:"required"`
		Creds  *credentials.Credentials
	}

	App struct {
		repo     *Repo
		usecases *Usecases
		api      *API
	}
)

func New(d *Dependencies) (*App, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.New", d, err)
	}

	repo, err := NewRepo(&RepoDependencies{})
	if err != nil {
		return nil, err
	}

	usecases, err := NewUsecases(&UsecasesDependencies{
		repo:   repo,
		config: d.Config,
	})
	if err != nil {
		return nil, err
	}

	api, err := NewAPI(&APIDependencies{
		config:   d.Config,
		usecases: usecases,
	})
	if err != nil {
		return nil, err
	}

	return &App{
		api:      api,
		repo:     repo,
		usecases: usecases,
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) <-chan error {
	errCh := make(chan error, 1)

	wg.Go(func() {
		if err := a.api.httpServer.Start(); err != nil {
			errCh <- fmt.Errorf("httpServer.Start error: %w", err)
		}
	})

	return errCh
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.api.httpServer.Shutdown(ctx); err != nil {
		slog.Error("http server shutdown failed", logger.Error(err))
	}
	slog.Debug("http server shutdown complete")
}
