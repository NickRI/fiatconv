package external

import (
	"context"

	"github.com/shopspring/decimal"
)

//go:generate mockgen -destination=../internal/mock/remote.client.go -package=mock github.com/NickRI/fiatconv/external RemoteClient
type RemoteClient interface {
	RequestRate(ctx context.Context, src, dst string) (decimal.Decimal, error)
}
