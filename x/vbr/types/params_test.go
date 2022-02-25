package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var validDistrEpochIdentifier = EpochDay

var validEarnRate = sdk.NewDecWithPrec(5, 1)
var invalidEarnRate = sdk.NewDecWithPrec(-5, 1)

func TestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  Params
		wantErr bool
	}{
		{
			name:   "ok",
			params: DefaultParams(),
		},
		{
			name:    "invalid DistrEpochIdentifier",
			params:  NewParams("", validEarnRate),
			wantErr: true,
		},
		{
			name: "invalid EarnRate",
			params: Params{
				DistrEpochIdentifier: validDistrEpochIdentifier,
				EarnRate:             invalidEarnRate,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := tt.params.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Params.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateDistrEpochIdentifier(t *testing.T) {
	tests := []struct {
		name                 string
		distrEpochIdentifier interface{}
		wantErr              bool
	}{
		{
			name:                 "ok",
			distrEpochIdentifier: validDistrEpochIdentifier,
		},
		{
			name:    "wrong type",
			wantErr: true,
		},
		{
			name:                 "empty",
			distrEpochIdentifier: "",
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateDistrEpochIdentifier(tt.distrEpochIdentifier); (err != nil) != tt.wantErr {
				t.Errorf("validateDistrEpochIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateEarnRate(t *testing.T) {
	tests := []struct {
		name     string
		earnRate interface{}
		wantErr  bool
	}{
		{
			name:     "ok",
			earnRate: validEarnRate,
		},
		{
			name:    "wrong type",
			wantErr: true,
		},
		{
			name:     "negative",
			earnRate: invalidEarnRate,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEarnRate(tt.earnRate); (err != nil) != tt.wantErr {
				t.Errorf("validateEarnRate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
