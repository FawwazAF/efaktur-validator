package efaktur

import (
	"context"

	"github.com/efaktur-validator/internal/model"
)

type djpRepository interface {
	GetInvoicesFromDJP(ctx context.Context, url string) (model.DJPEfaktur, error)
}

type efakturController struct {
	djp djpRepository
}

func New(djp djpRepository) *efakturController {
	return &efakturController{
		djp: djp,
	}
}
