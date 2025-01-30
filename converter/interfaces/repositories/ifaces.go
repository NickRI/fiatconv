package repositories

import (
	"context"

	"github.com/NickRI/fiatconv/converter/domain/entries"
)

//go:generate mockgen -destination=../../../internal/mock/repository.go -package=mock github.com/NickRI/fiatconv/converter/interfaces/repositories Repository
type Repository interface {
	GetItemPrice(ctx context.Context, src, dst entries.CurrencySymbol) (entries.Amount, error)
}
