package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

const apiAddress = "https://api.exchangeratesapi.io/latest"

type errorResponse struct {
	Err string `json:"error"`
}

type ratesResponse struct {
	Rates map[string]decimal.Decimal `json:"rates"`
	Base  string                     `json:"base"`
	Date  string                     `json:"date"`
}

type ExchangeRatesClient struct {
	base *CtxHttpClient
}

func NewExchangeRatesClient() RemoteClient {
	return &ExchangeRatesClient{base: NewCtxHttpClient()}
}

func (e *ExchangeRatesClient) RequestRate(ctx context.Context, src, dst string) (decimal.Decimal, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?symbols=%s&base=%s", apiAddress, dst, src), nil)
	if err != nil {
		return decimal.Zero, xerrors.Errorf("exchangerates: error during creating new request: %w", err)
	}

	resp, err := e.base.Do(ctx, req.WithContext(ctx))
	if err != nil {
		return decimal.Zero, xerrors.Errorf("exchangerates: error during making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errResp := errorResponse{}
		if err = json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return decimal.Zero, xerrors.Errorf("exchangerates: error during decoding error response: %w", err)
		}
		return decimal.Zero, xerrors.Errorf("exchangerates: error during converting restapi: %s", errResp.Err)
	}

	ratesResp := ratesResponse{}
	if err = json.NewDecoder(resp.Body).Decode(&ratesResp); err != nil {
		return decimal.Zero, xerrors.Errorf("exchangerates: error during decoding response: %w", err)
	}

	rate, ok := ratesResp.Rates[dst]
	if !ok {
		return rate, xerrors.New("exchangerates: distenation error not found")
	}

	return rate, nil
}
