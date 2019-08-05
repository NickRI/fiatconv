package entries

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

type CurrencySymbol string

func (cs CurrencySymbol) String() string {
	return string(cs)
}

type CurrencySymbols [2]CurrencySymbol

func (css CurrencySymbols) String() string {
	return fmt.Sprintf("%s:%s", css[0], css[1])
}

type Amount struct {
	decimal.Decimal
}

func (a *Amount) Mul(b Amount) Amount {
	return Amount{a.Decimal.Mul(b.Decimal)}
}

var Zero = Amount{decimal.Zero}

var KeyNotFoundError = xerrors.New("key not found")

type CachedAmount struct {
	Until  time.Time
	Amount Amount
}

func NewCachedAmount(until time.Time, amount Amount) CachedAmount {
	return CachedAmount{Until: until, Amount: amount}
}

func (ca *CachedAmount) Before(t time.Time) bool {
	return ca.Until.Before(t)
}
