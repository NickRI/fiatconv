package entries

import "github.com/shopspring/decimal"

type CurrencySymbol string

func (cs CurrencySymbol) String() string {
	return string(cs)
}

type Amount struct {
	decimal.Decimal
}

func (a *Amount) Mul(b Amount) Amount {
	return Amount{a.Decimal.Mul(b.Decimal)}
}

var Zero = Amount{decimal.Zero}
