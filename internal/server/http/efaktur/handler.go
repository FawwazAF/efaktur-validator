package efaktur

import (
	"context"

	"github.com/efaktur-validator/internal/model"
)

type efakturControllerAdapter interface {
	ValidateEfaktur(ctx context.Context, pdfPath string) (model.EfakturValidationResult, error)
}

type Handler struct {
	efaktur efakturControllerAdapter
}

func New(efaktur efakturControllerAdapter) *Handler {
	return &Handler{
		efaktur: efaktur,
	}
}
