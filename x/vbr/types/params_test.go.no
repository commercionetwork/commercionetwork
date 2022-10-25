package types

import (
	"testing"
)

func TestParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  Params
		wantErr bool
	}{
		{
			name:   "ok",
			params: validParams,
		},
		{
			name:    "invalid DistrEpochIdentifier",
			params:  NewParams("", validEarnRate),
			wantErr: true,
		},
		{
			name:    "invalid EarnRate",
			params:  NewParams(validDistrEpochIdentifier, InvalidEarnRate),
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
			if err := validateDistrEpochIdentifierParamSetPairs(tt.distrEpochIdentifier); (err != nil) != tt.wantErr {
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
			earnRate: InvalidEarnRate,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEarnRateParamSetPairs(tt.earnRate); (err != nil) != tt.wantErr {
				t.Errorf("validateEarnRate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
