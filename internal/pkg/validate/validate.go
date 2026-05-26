package validate

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validationErrors *validator.Validate
	initOnce         sync.Once
)

func Struct(obj any) error {
	return get().Struct(obj)
}

func get() *validator.Validate {
	initOnce.Do(func() {
		validationErrors = validator.New()
	})

	return validationErrors
}
