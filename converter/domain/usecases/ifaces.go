package usecases

import (
	"context"

	"github.com/NickRI/fiatconv/converter/domain/entries"
)

//go:generate mockgen -destination=../../../internal/mock/usecase.go -package=mock github.com/NickRI/fiatconv/converter/domain/usecases Usecase
type Usecase interface {
	Convert(ctx context.Context, src, dst entries.CurrencySymbol, amount entries.Amount) (entries.Amount, error)
}
