package external

import (
	"context"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/shopspring/decimal"
)

//go:generate mockgen -destination=../internal/mock/remote.client.go -package=mock github.com/NickRI/fiatconv/external RemoteClient
type RemoteClient interface {
	RequestRate(ctx context.Context, src, dst string) (decimal.Decimal, error)
}

type KeyValStorage interface {
	Get(key string) (entries.CachedAmount, error)
	Set(key string, val entries.CachedAmount) error
}
