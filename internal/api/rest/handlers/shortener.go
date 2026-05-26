package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/SaTeR151/url_shortener/internal/pkg/errorspkg"
	"github.com/SaTeR151/url_shortener/internal/pkg/logger"
	"github.com/SaTeR151/url_shortener/internal/pkg/validate"
)

type ShortenerManager interface {
	Create(context.Context, string) (string, error)
	Get(context.Context, string) (string, error)
}

type (
	ShortenerDependencies struct {
		ShortenerManager ShortenerManager `validate:"required"`
	}

	ShortenerHandlers struct {
		shortenerManager ShortenerManager
		logger           *slog.Logger
	}
)

func NewShortenerHandler(d *ShortenerDependencies) (*ShortenerHandlers, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("handlers.NewShortenerHandler", d, err)
	}

	return &ShortenerHandlers{
		shortenerManager: d.ShortenerManager,
		logger:           logger.WithComponent("handlers.Shortener"),
	}, nil
}

func (h *ShortenerHandlers) Create(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		h.sendError(res, errors.New("method not allowed"))
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		err = errorspkg.NewReadBodyError(err)
		h.sendError(res, err)
		return
	}

	shortenedURL, err := h.shortenerManager.Create(req.Context(), string(body))
	if err != nil {
		h.sendError(res, err)
		return
	}

	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte(shortenedURL))
	res.WriteHeader(http.StatusCreated)
}

func (h *ShortenerHandlers) Get(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		h.sendError(res, errors.New("method not allowed"))
		return
	}

	code := req.PathValue("code")
	if code == "" {
		h.sendError(res, errorspkg.NewIsEmptyError("code"))
		return
	}

	defaultURL, err := h.shortenerManager.Get(req.Context(), code)
	if err != nil {
		h.sendError(res, err)
		return
	}

	http.Redirect(res, req, defaultURL, http.StatusMovedPermanently)
}

func (h *ShortenerHandlers) sendError(res http.ResponseWriter, err error) {
	message, status := h.parseError(err)
	http.Error(res, message, status)
}

func (h *ShortenerHandlers) parseError(err error) (string, int) {
	h.handleError(err)

	switch {
	case errors.As(err, &errorspkg.IsEmptyError{}),
		errors.As(err, &errorspkg.NotValidError{}):
		return err.Error(), http.StatusBadRequest
	case errors.As(err, &errorspkg.NotFoundError{}):
		return err.Error(), http.StatusNotFound
	case errors.As(err, &errorspkg.FailedError{}):
		return err.Error(), http.StatusInternalServerError
	default:
		return errors.New("internal server error").Error(), http.StatusInternalServerError
	}
}

func (h *ShortenerHandlers) handleError(err error) {
	h.logger.Error(err.Error())
}
