package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateConversionRate(t *testing.T) {

	tests := []struct {
		name    string
		rate    sdk.Dec
		wantErr bool
	}{
		{
			name:    "ok",
			rate:    sdk.NewDec(1),
			wantErr: false,
		},
		{
			name:    "zero",
			rate:    sdk.NewDec(0),
			wantErr: true,
		},
		{
			name:    "negative",
			rate:    sdk.NewDec(-1),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConversionRate(tt.rate); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConversionRate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFreezePeriod(t *testing.T) {
	type args struct {
		freezePeriod time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok positive",
			args: args{
				freezePeriod: time.Minute,
			},
			wantErr: false,
		},
		{
			name: "ok zero",
			args: args{
				freezePeriod: 0,
			},
			wantErr: false,
		},
		{
			name: "negative",
			args: args{
				freezePeriod: -time.Minute,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateFreezePeriod(tt.args.freezePeriod); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFreezePeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
