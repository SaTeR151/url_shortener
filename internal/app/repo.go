package app

import (
	"github.com/SaTeR151/url_shortener/internal/pkg/errorspkg"
	"github.com/SaTeR151/url_shortener/internal/pkg/validate"
	inmemory "github.com/SaTeR151/url_shortener/internal/repo/inMemory"
)

type (
	RepoDependencies struct{}

	Repo struct {
		inMemory *inmemory.ShortenerRepo
	}
)

func NewRepo(d *RepoDependencies) (*Repo, error) {
	if err := validate.Struct(d); err != nil {
		return nil, errorspkg.NewValidationError("app.NewRepo", d, err)
	}

	inMemory, err := inmemory.New(&inmemory.Dependencies{})
	if err != nil {
		return nil, err
	}

	return &Repo{
		inMemory: inMemory,
	}, nil
}
