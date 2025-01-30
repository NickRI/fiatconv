package repositories

import (
	"context"
	"time"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/NickRI/fiatconv/external"
	"golang.org/x/xerrors"
)

type CachedRepo struct {
	Repository
	cache      external.KeyValStorage
	expiration time.Duration
}

func NewCachedRepo(r Repository, k external.KeyValStorage, e time.Duration) Repository {
	return &CachedRepo{Repository: r, cache: k, expiration: e}
}

func (c *CachedRepo) GetItemPrice(ctx context.Context, src, dst entries.CurrencySymbol) (entries.Amount, error) {
	until := time.Now().Add(c.expiration)
	symbols := &entries.CurrencySymbols{src, dst}

	cached, err := c.cache.Get(symbols.String())

	if err != nil {
		if !xerrors.Is(err, entries.KeyNotFoundError) {
			return entries.Zero, err
		}
	}

	if err != nil || cached.Before(until) {

		tpe, err := c.Repository.GetItemPrice(ctx, src, dst)
		if err != nil {
			return entries.Zero, err
		}

		camo := entries.NewCachedAmount(until, tpe)

		err = c.cache.Set(symbols.String(), camo)
		if err != nil {
			return entries.Zero, err
		}
		return tpe, nil
	}

	return cached.Amount, nil
}
