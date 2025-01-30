package presenters

import (
	"context"
	"io"

	"github.com/NickRI/fiatconv/converter/domain/entries"
)

//go:generate mockgen -destination=../../../internal/mock/presenter.go -package=mock github.com/NickRI/fiatconv/converter/interfaces/presenters Presenter
type Presenter interface {
	Present(ctx context.Context, src, dst entries.CurrencySymbol, in, conv entries.Amount) error
	Error(ctx context.Context, err error) error
}

//StringWriter local wrapper to generate mocks
//go:generate mockgen -destination=../../../internal/mock/string.writer.go -package=mock github.com/NickRI/fiatconv/converter/interfaces/presenters StringWriter
type StringWriter interface {
	io.StringWriter
}
