package controllers

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/shopspring/decimal"

	"golang.org/x/xerrors"

	"github.com/NickRI/fiatconv/internal/mock"

	"github.com/golang/mock/gomock"
)

type xError struct {
	err error
}

func XError(e error) gomock.Matcher {
	return &xError{e}
}

func (xe *xError) Matches(x interface{}) bool {
	return reflect.TypeOf(x).String() == reflect.TypeOf(xe.err).String() &&
		fmt.Sprintf("%s", x) == xe.err.Error()
}

func (xe *xError) String() string {
	return reflect.TypeOf(xe.err).String() + " " + xe.err.Error()
}

func TestCli_Convert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInterator := mock.NewMockUsecase(ctrl)
	mockPresenter := mock.NewMockPresenter(ctrl)

	testError := xerrors.New("some_error")
	testAmount := entries.Amount{decimal.NewFromFloat(4.71)}

	type args struct {
		ctx  context.Context
		args []string
	}
	tests := []struct {
		name    string
		before  func(*args)
		args    args
		wantErr error
	}{
		{
			name: "wrong arguments count",
			before: func(a *args) {
				mockPresenter.EXPECT().Error(a.ctx, XError(xerrors.Errorf("wrong number of arguments: %d", 2))).Return(nil)
			},
			args: args{
				ctx:  context.Background(),
				args: []string{"", ""},
			},
		},
		{
			name: "wrong amount",
			before: func(a *args) {
				err := xerrors.New("can't convert wrong_amount to decimal")
				mockPresenter.EXPECT().Error(a.ctx, XError(xerrors.Errorf("amount conversion err: %w", err))).Return(nil)
			},
			args: args{
				ctx:  context.Background(),
				args: []string{"wrong_amount", "USD", "RUB"},
			},
		},
		{
			name: "repository return error",
			before: func(a *args) {
				dam, _ := decimal.NewFromString(a.args[0])
				amount := entries.Amount{dam}
				src := entries.CurrencySymbol(strings.ToUpper(a.args[1]))
				dst := entries.CurrencySymbol(strings.ToUpper(a.args[2]))

				mockInterator.EXPECT().Convert(a.ctx, src, dst, amount).Return(entries.Zero, testError)
				mockPresenter.EXPECT().Error(a.ctx, XError(xerrors.Errorf("error during processing conversion: %w", testError))).Return(nil)
			},
			args: args{
				ctx:  context.Background(),
				args: []string{"3.53", "USD", "RUB"},
			},
		},
		{
			name: "works well",
			before: func(a *args) {
				dam, _ := decimal.NewFromString(a.args[0])
				amount := entries.Amount{dam}
				src := entries.CurrencySymbol(strings.ToUpper(a.args[1]))
				dst := entries.CurrencySymbol(strings.ToUpper(a.args[2]))

				mockInterator.EXPECT().Convert(a.ctx, src, dst, amount).Return(testAmount, nil)
				mockPresenter.EXPECT().Present(a.ctx, src, dst, amount, testAmount).Return(nil)
			},
			args: args{
				ctx:  context.Background(),
				args: []string{"3.53", "USD", "RUB"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCli(mockInterator, mockPresenter)
			tt.before(&tt.args)
			err := c.Convert(tt.args.ctx, tt.args.args)
			if err != nil && err.Error() != tt.wantErr.Error() || tt.wantErr != nil && err == nil {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
