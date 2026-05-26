package rest

import (
	"net/http"

	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/validate"
)

type (
	ShortenerHandlerRouter interface {
		Create(http.ResponseWriter, *http.Request)
		Get(http.ResponseWriter, *http.Request)
	}
)

type (
	RouterDependencies struct {
		ShortenerHandlers ShortenerHandlerRouter `validate:"required"`
	}
)

func NewRouter(d *RouterDependencies) (http.Handler, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewRouter", d, err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/{code}", d.ShortenerHandlers.Get)
	mux.HandleFunc("/", d.ShortenerHandlers.Create)

	return mux, nil
}
