package shortener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"

	"github.com/sater151/url_shortener/internal/config"
	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/generator"
	"github.com/sater151/url_shortener/internal/pkg/logger"
	"github.com/sater151/url_shortener/internal/pkg/validate"
)

type (
	Repo interface {
		Create(_ context.Context, defaultURL, shortURL string) error
		Get(_ context.Context, shortURL string) (string, error)
		IsURLShorted(_ context.Context, defaultURL string) (string, error)
	}
)

type (
	Dependencies struct {
		Repo   Repo              `validate:"required"`
		Config *config.Shortener `validate:"required"`
	}

	Manager struct {
		repo            Repo
		config          *config.Shortener
		addressTemplate string
		logger          *slog.Logger
	}
)

func New(d *Dependencies) (*Manager, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("shortener.New", d, err)
	}

	m := &Manager{
		repo:   d.Repo,
		config: d.Config,
		logger: logger.WithComponent("usecases.Shortener"),
	}

	m.buildURL()

	return m, nil
}

func (m *Manager) Create(ctx context.Context, url string) (string, error) {
	if err := m.isURL(url); err != nil {
		return "", err
	}

	code, err := m.repo.IsURLShorted(ctx, url)
	if err != nil && errors.As(err, &errorspkg.NotFoundError{}) {
		return "", err
	}

	if code != "" {
		return m.codeInURL(code), nil
	}

	// Цикл, чтобы в случае создания кода, по которому уже есть ссылка - создавали новый.
	for i := 0; i < m.config.CountOfRetriesCreatingLink; i++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			code := generator.Code(m.config.CountOfSymbolsInURL)

			_, err := m.repo.Get(ctx, code)
			if err != nil {
				if errors.As(err, &errorspkg.NotFoundError{}) {
					err = m.repo.Create(ctx, url, code)
					if err != nil {
						return "", err
					}

					return m.codeInURL(code), nil
				}

				return "", err
			}

			// На случай, если во время сохранения пользователем URL кто-то другой создаст на эту же URL сокращенную ссылку.
			code, err = m.repo.IsURLShorted(ctx, url)
			if err != nil && errors.As(err, &errorspkg.NotFoundError{}) {
				return "", err
			}

			if code != "" {
				return m.codeInURL(code), nil
			}
		}
	}

	return "", errorspkg.NewFailedError("create link")
}

func (m *Manager) Get(ctx context.Context, code string) (string, error) {
	return m.repo.Get(ctx, code)
}

func (m *Manager) codeInURL(code string) string {
	return fmt.Sprint(m.addressTemplate, code)
}

func (m *Manager) buildURL() {
	schema := "http"
	if m.config.Secure {
		schema = "https"
	}

	m.addressTemplate = fmt.Sprintf("%s://%s/%s", schema, m.config.Host, addressTemplate)
}

func (m *Manager) isURL(checkedURL string) error {
	if checkedURL == "" {
		return errorspkg.NewIsEmptyError("url")
	}

	_, err := url.ParseRequestURI(checkedURL)
	if err != nil {
		return errorspkg.NewNotValidError("text", "url")
	}

	return nil
}

func (m *Manager) handleError(err error) {
	m.logger.Error(err.Error())
}
