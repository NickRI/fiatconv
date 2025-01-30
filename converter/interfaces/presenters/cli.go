package presenters

import (
	"context"
	"fmt"

	"github.com/NickRI/fiatconv/converter/domain/entries"
)

const usage = `Fiatconv utility converting inputted currencies between each other using https://exchangeratesapi.io
arguments: <amount:float> <src_symbol:string> <dst_symbol:string>
example: fiatconv 123.45 USD RUB`

type Cli struct {
	writer StringWriter
}

func NewCli(w StringWriter) Presenter {
	return &Cli{writer: w}
}

func (c *Cli) Present(ctx context.Context, src, dst entries.CurrencySymbol, in, conv entries.Amount) error {
	_, err := c.writer.WriteString(fmt.Sprintf("%s %s = %s %s\n", in.StringFixedBank(2), src, conv.StringFixedBank(2), dst))
	return err
}

func (c *Cli) Error(ctx context.Context, err error) error {
	if _, err := c.writer.WriteString(usage); err != nil {
		return err
	}
	if _, err := c.writer.WriteString(fmt.Sprintf("\n\n%s\n", err.Error())); err != nil {
		return err
	}
	return nil
}
