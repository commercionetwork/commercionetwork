package types

import (
	"testing"
	"time"

	"cosmossdk.io/math"
)

func TestValidateConversionRate(t *testing.T) {

	tests := []struct {
		name    string
		rate    math.LegacyDec
		wantErr bool
	}{
		{
			name:    "ok",
			rate:    math.LegacyNewDec(1),
			wantErr: false,
		},
		{
			name:    "zero",
			rate:    math.LegacyNewDec(0),
			wantErr: true,
		},
		{
			name:    "negative",
			rate:    math.LegacyNewDec(-1),
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

var validConversionRate = math.LegacyNewDec(2)
var invalidConversionRate = math.LegacyNewDec(-1)
var validFreezePeriod time.Duration = 0
var invalidFreezePeriod = -time.Minute

// var validParams = Params{
// 	ConversionRate: validConversionRate,
// 	FreezePeriod:   validFreezePeriod,
// }

func TestParams_Validate(t *testing.T) {
	type fields struct {
		ConversionRate math.LegacyDec
		FreezePeriod   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				ConversionRate: validConversionRate,
				FreezePeriod:   validFreezePeriod,
			},
		},
		{
			name: "invalid conversion rate",
			fields: fields{
				ConversionRate: invalidConversionRate,
				FreezePeriod:   validFreezePeriod,
			},
			wantErr: true,
		},
		{
			name: "invalid freeze period",
			fields: fields{
				ConversionRate: validConversionRate,
				FreezePeriod:   invalidFreezePeriod,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				ConversionRate: tt.fields.ConversionRate,
				FreezePeriod:   tt.fields.FreezePeriod,
			}
			if err := p.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Params.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateConversionRateParamSetPairs(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				i: validConversionRate,
			},
		},
		{
			name: "invalid type",
			args: args{
				i: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConversionRateParamSetPairs(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("validateConversionRateParamSetPairs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateFreezePeriodParamSetPairs(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				i: validFreezePeriod,
			},
		},
		{
			name:    "invalid type",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFreezePeriodParamSetPairs(tt.args.i); (err != nil) != tt.wantErr {
				t.Errorf("validateFreezePeriodParamSetPairs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
