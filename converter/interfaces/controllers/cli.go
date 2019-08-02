package controllers

import (
	"context"
	"strings"

	"github.com/shopspring/decimal"

	"github.com/NickRI/fiatconv/converter/domain/entries"

	"golang.org/x/xerrors"

	"github.com/NickRI/fiatconv/converter/domain/usecases"
	"github.com/NickRI/fiatconv/converter/interfaces/presenters"
)

type Cli struct {
	interactor usecases.Usecase
	presenter  presenters.Presenter
}

func NewCli(i usecases.Usecase, p presenters.Presenter) *Cli {
	return &Cli{
		interactor: i,
		presenter:  p,
	}
}

func (c *Cli) Convert(ctx context.Context, args []string) error {
	if len(args) < 3 {
		return c.presenter.Error(ctx, xerrors.Errorf("wrong number of arguments: %d", len(args)))
	}

	dam, err := decimal.NewFromString(args[0])
	if err != nil {
		return c.presenter.Error(ctx, xerrors.Errorf("amount conversion err: %w", err))
	}

	amount := entries.Amount{dam}
	src := entries.CurrencySymbol(strings.ToUpper(args[1]))
	dst := entries.CurrencySymbol(strings.ToUpper(args[2]))

	converted, err := c.interactor.Convert(ctx, src, dst, amount)
	if err != nil {
		return c.presenter.Error(ctx, xerrors.Errorf("error during processing conversion: %w", err))
	}
	return c.presenter.Present(ctx, src, dst, amount, converted)
}
