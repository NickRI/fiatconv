package presenters

import (
	"context"
	"fmt"
	"testing"

	"github.com/shopspring/decimal"

	"golang.org/x/xerrors"

	"github.com/NickRI/fiatconv/converter/domain/entries"
	"github.com/NickRI/fiatconv/internal/mock"
	"github.com/golang/mock/gomock"
)

func TestCli_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	writerMock := mock.NewMockStringWriter(ctrl)
	testError := xerrors.New("some_error")
	testError2 := xerrors.New("some_error_2")

	type args struct {
		ctx context.Context
		err error
	}
	tests := []struct {
		name    string
		before  func(*args)
		args    args
		wantErr error
	}{
		{
			name: "usage writer returns error",
			before: func(a *args) {
				writerMock.EXPECT().WriteString(usage).Return(0, testError)
			},
			args: args{
				ctx: context.Background(),
				err: xerrors.New("test"),
			},
			wantErr: testError,
		},
		{
			name: "message writer returns error",
			before: func(a *args) {
				writerMock.EXPECT().WriteString(usage).Return(len(usage), nil)
				writerMock.EXPECT().WriteString(fmt.Sprintf("\n\n%s\n", a.err.Error())).Return(0, testError2)
			},
			args: args{
				ctx: context.Background(),
				err: xerrors.New("test_2"),
			},
			wantErr: testError2,
		},
		{
			name: "working well",
			before: func(a *args) {
				errString := fmt.Sprintf("\n\n%s\n", a.err.Error())
				writerMock.EXPECT().WriteString(usage).Return(len(usage), nil)
				writerMock.EXPECT().WriteString(errString).Return(len(errString), nil)
			},
			args: args{
				ctx: context.Background(),
				err: xerrors.New("test_3"),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCli(writerMock)
			tt.before(&tt.args)
			err := c.Error(tt.args.ctx, tt.args.err)
			if err != nil && !xerrors.Is(err, tt.wantErr) || tt.wantErr != nil && err == nil {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCli_Present(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	writerMock := mock.NewMockStringWriter(ctrl)
	testError := xerrors.New("some_error")
	testIn := entries.Amount{decimal.NewFromFloat(3.5343434)}
	testOut := entries.Amount{decimal.NewFromFloat(42.1)}

	type args struct {
		ctx  context.Context
		src  entries.CurrencySymbol
		dst  entries.CurrencySymbol
		in   entries.Amount
		conv entries.Amount
	}
	tests := []struct {
		name    string
		before  func(*args)
		args    args
		wantErr error
	}{
		{
			name: "writer returns error",
			before: func(a *args) {
				output := fmt.Sprintf("%s %s = %s %s\n", a.in.StringFixedBank(2), a.src, a.conv.StringFixedBank(2), a.dst)
				writerMock.EXPECT().WriteString(output).Return(0, testError)
			},
			args: args{
				ctx:  context.Background(),
				src:  "USD",
				dst:  "RUB",
				in:   testIn,
				conv: testOut,
			},
			wantErr: testError,
		},
		{
			name: "working well",
			before: func(a *args) {
				output := fmt.Sprintf("%s %s = %s %s\n", a.in.StringFixedBank(2), a.src, a.conv.StringFixedBank(2), a.dst)
				writerMock.EXPECT().WriteString(output).Return(len(output), nil)
			},
			args: args{
				ctx:  context.Background(),
				src:  "USD",
				dst:  "RUB",
				in:   testIn,
				conv: testOut,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCli(writerMock)
			tt.before(&tt.args)
			err := c.Present(tt.args.ctx, tt.args.src, tt.args.dst, tt.args.in, tt.args.conv)
			if err != nil && !xerrors.Is(err, tt.wantErr) || tt.wantErr != nil && err == nil {
				t.Errorf("Present() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
