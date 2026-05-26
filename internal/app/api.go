package app

import (
	"github.com/sater151/url_shortener/internal/api/rest"
	"github.com/sater151/url_shortener/internal/api/rest/handlers"
	"github.com/sater151/url_shortener/internal/config"
	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/validate"
)

type (
	APIDependencies struct {
		config   *config.Config
		usecases *Usecases
	}

	API struct {
		httpServer *rest.HTTPServer
	}
)

func NewAPI(d *APIDependencies) (*API, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewAPI", d, err)
	}

	shortenerHandler, err := handlers.NewShortenerHandler(&handlers.ShortenerDependencies{
		ShortenerManager: d.usecases.shortener,
	})
	if err != nil {
		return nil, err
	}

	router, err := rest.NewRouter(&rest.RouterDependencies{
		ShortenerHandlers: shortenerHandler,
	})
	if err != nil {
		return nil, err
	}

	server, err := rest.NewHTTPServer(&rest.HTTPServerDependencies{
		Handler: router,
		Config:  d.config.HTTPServer,
	})
	if err != nil {
		return nil, err
	}

	return &API{
		httpServer: server,
	}, nil
}
