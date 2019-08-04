package usecases

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"

	"golang.org/x/xerrors"

	"github.com/NickRI/fiatconv/internal/mock"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/golang/mock/gomock"
)

func TestBaseInteractor_Convert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoMock := mock.NewMockRepository(ctrl)
	testAmount := entries.Amount{decimal.NewFromFloat32(7.42)}
	testError := xerrors.New("some_error")
	type args struct {
		ctx    context.Context
		src    entries.CurrencySymbol
		dst    entries.CurrencySymbol
		amount entries.Amount
	}
	tests := []struct {
		name    string
		args    args
		before  func(*args)
		want    entries.Amount
		wantErr error
	}{
		{
			name: "repository error",
			args: args{
				ctx:    context.Background(),
				src:    "RUB",
				dst:    "EUR",
				amount: entries.Zero,
			},
			before: func(a *args) {
				repoMock.EXPECT().GetItemPrice(a.ctx, a.src, a.dst).Return(entries.Zero, testError)
			},
			want:    entries.Zero,
			wantErr: testError,
		},
		{
			name: "passing zero",
			args: args{
				ctx:    context.Background(),
				src:    "RUB",
				dst:    "EUR",
				amount: entries.Zero,
			},
			before: func(a *args) {
				repoMock.EXPECT().GetItemPrice(a.ctx, a.src, a.dst).Return(testAmount, nil)
			},
			want: entries.Zero,
		},
		{
			name: "service returns zero",
			args: args{
				ctx:    context.Background(),
				src:    "RUB",
				dst:    "EUR",
				amount: testAmount,
			},
			before: func(a *args) {
				repoMock.EXPECT().GetItemPrice(a.ctx, a.src, a.dst).Return(entries.Zero, nil)
			},
			want: entries.Zero,
		},
		{
			name: "passing zero & service returns zero",
			args: args{
				ctx:    context.Background(),
				src:    "RUB",
				dst:    "EUR",
				amount: entries.Zero,
			},
			before: func(a *args) {
				repoMock.EXPECT().GetItemPrice(a.ctx, a.src, a.dst).Return(entries.Zero, nil)
			},
			want: entries.Zero,
		},
		{
			name: "normal multiply",
			args: args{
				ctx:    context.Background(),
				src:    "RUB",
				dst:    "EUR",
				amount: testAmount,
			},
			before: func(a *args) {
				repoMock.EXPECT().GetItemPrice(a.ctx, a.src, a.dst).Return(testAmount, nil)
			},
			want: testAmount.Mul(testAmount),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := NewBaseInteractor(repoMock)
			tt.before(&tt.args)
			got, err := ex.Convert(tt.args.ctx, tt.args.src, tt.args.dst, tt.args.amount)
			if err != nil && !xerrors.Is(err, tt.wantErr) || tt.wantErr != nil && err == nil {
				t.Errorf("Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want.Decimal) {
				t.Errorf("Convert() got = %v, want %v", got, tt.want)
			}
		})
	}
}
