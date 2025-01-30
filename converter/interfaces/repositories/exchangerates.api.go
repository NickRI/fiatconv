package repositories

import (
	"context"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"golang.org/x/xerrors"

	"github.com/NickRI/fiatconv/external"
)

type RestapiRepo struct {
	client external.RemoteClient
}

func NewRestapiRepo(c external.RemoteClient) Repository {
	return &RestapiRepo{
		client: c,
	}
}

func (r *RestapiRepo) GetItemPrice(ctx context.Context, src, dst entries.CurrencySymbol) (entries.Amount, error) {
	itp, err := r.client.RequestRate(ctx, src.String(), dst.String())
	if err != nil {
		return entries.Zero, xerrors.Errorf("error getting response: %w", err)
	}

	return entries.Amount{itp}, err
}
