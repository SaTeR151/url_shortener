package inmemory

import (
	"context"
	"log/slog"

	"github.com/sater151/url_shortener/internal/pkg/errorspkg"
	"github.com/sater151/url_shortener/internal/pkg/logger"
	"github.com/sater151/url_shortener/internal/pkg/validate"
)

type (
	Dependencies struct{}

	ShortenerRepo struct {
		logger    *slog.Logger
		existURLs map[defaultURL]code
		shortURLs map[code]defaultURL
	}
)

func New(d *Dependencies) (*ShortenerRepo, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("inmemory.New", d, err)
	}

	return &ShortenerRepo{
		logger: logger.WithComponent("repo.InMemory"),
	}, nil
}

func (m *ShortenerRepo) Create(_ context.Context, defaultURL, shortURL string) error {
	m.shortURLs[inCode(shortURL)] = inDefaultURL(defaultURL)
	m.existURLs[inDefaultURL(defaultURL)] = inCode(shortURL)

	return nil
}

func (m *ShortenerRepo) Get(_ context.Context, shortURL string) (string, error) {
	if defaultURL, ok := m.shortURLs[inCode(shortURL)]; ok {
		return defaultURL.String(), nil
	}

	return "", errorspkg.NewNotFoundError("URL", shortURL)
}

func (m *ShortenerRepo) IsURLShorted(_ context.Context, defaultURL string) (string, error) {
	if shortURL, ok := m.existURLs[inDefaultURL(defaultURL)]; ok {
		return shortURL.String(), nil
	}

	return "", errorspkg.NewNotFoundError("shorted URL", defaultURL)
}
