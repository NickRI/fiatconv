package usecases

import (
	"context"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/NickRI/fiatconv/converter/interfaces/repositories"
	"golang.org/x/xerrors"
)

type BaseInteractor struct {
	repo repositories.Repository
}

func NewBaseInteractor(r repositories.Repository) Usecase {
	return &BaseInteractor{
		repo: r,
	}
}

func (ex *BaseInteractor) Convert(ctx context.Context, src, dst entries.CurrencySymbol, amount entries.Amount) (entries.Amount, error) {
	itp, err := ex.repo.GetItemPrice(ctx, src, dst)
	if err != nil {
		return itp, xerrors.Errorf("error from repository: %w", err)
	}

	return itp.Mul(amount), nil
}
