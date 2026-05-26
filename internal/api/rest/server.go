package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/SaTeR151/url_shortener/internal/config"
	"github.com/SaTeR151/url_shortener/internal/pkg/errorspkg"
	"github.com/SaTeR151/url_shortener/internal/pkg/validate"
)

type (
	HTTPServerDependencies struct {
		Handler http.Handler       `validate:"required"`
		Config  *config.HTTPServer `validate:"required"`
	}

	HTTPServer struct {
		srvHTTP *http.Server
	}
)

func NewHTTPServer(d *HTTPServerDependencies) (*HTTPServer, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewHTTPServer", d, err)
	}

	return &HTTPServer{
		srvHTTP: &http.Server{
			Addr:              ":" + d.Config.Port,
			Handler:           d.Handler,
			ReadTimeout:       d.Config.ReadTimeout,
			ReadHeaderTimeout: d.Config.ReadHeaderTimeout,
			WriteTimeout:      d.Config.WriteTimeout,
			IdleTimeout:       d.Config.IdleTimeout,
			MaxHeaderBytes:    d.Config.MaxHeaderBytes,
		},
	}, nil
}

func (s *HTTPServer) Start() error {
	if err := s.srvHTTP.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.srvHTTP.Shutdown(ctx)
}
