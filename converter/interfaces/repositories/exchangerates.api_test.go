package repositories

import (
	"context"
	"testing"

	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/NickRI/fiatconv/internal/mock"
	"github.com/golang/mock/gomock"
)

func TestRestapiRepo_GetItemPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mock.NewMockRemoteClient(ctrl)
	testAmount := entries.Amount{decimal.NewFromFloat32(7.42)}
	testError := xerrors.New("some_error")

	type args struct {
		ctx context.Context
		src entries.CurrencySymbol
		dst entries.CurrencySymbol
	}
	tests := []struct {
		name    string
		args    args
		before  func(*args)
		want    entries.Amount
		wantErr error
	}{
		{
			name: "repository returns error",
			before: func(a *args) {
				clientMock.EXPECT().RequestRate(a.ctx, a.src.String(), a.dst.String()).Return(testAmount.Decimal, testError)
			},
			args: args{
				ctx: context.Background(),
				src: "RUB",
				dst: "USD",
			},
			want:    entries.Zero,
			wantErr: testError,
		},
		{
			name: "no errors",
			before: func(a *args) {
				clientMock.EXPECT().RequestRate(a.ctx, a.src.String(), a.dst.String()).Return(testAmount.Decimal, nil)
			},
			args: args{
				ctx: context.Background(),
				src: "RUB",
				dst: "USD",
			},
			want: testAmount,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRestapiRepo(clientMock)
			tt.before(&tt.args)
			got, err := r.GetItemPrice(tt.args.ctx, tt.args.src, tt.args.dst)
			if err != nil && !xerrors.Is(err, tt.wantErr) || tt.wantErr != nil && err == nil {
				t.Errorf("GetItemPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Equal(tt.want.Decimal) {
				t.Errorf("GetItemPrice() got = %v, want %v", got, tt.want)
			}
		})
	}
}
