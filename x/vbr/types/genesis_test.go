package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var validPoolAmount = sdk.NewDecCoins(sdk.NewDecCoin(BondDenom, sdk.NewInt(100)))

var validGenesis = GenesisState{
	PoolAmount: validPoolAmount,
	Params:     validParams,
}

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		PoolAmount sdk.DecCoins
		Params     Params
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields(validGenesis),
		},
		{
			name: "nil pool amount",
			fields: fields{
				PoolAmount: nil,
				Params:     validParams,
			},
			wantErr: true,
		},
		{
			name: "empty PoolAmount",
			fields: fields{
				PoolAmount: sdk.NewDecCoins(),
				Params:     validParams,
			},
			wantErr: true,
		},
		{
			name: "invalid DistrEpochIdentifier",
			fields: fields{
				PoolAmount: validPoolAmount,
				Params:     NewParams("", validEarnRate),
			},
			wantErr: true,
		},
		{
			name: "invalid EarnRate",
			fields: fields{
				PoolAmount: validPoolAmount,
				Params:     NewParams(validDistrEpochIdentifier, InvalidEarnRate),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := GenesisState{
				PoolAmount: tt.fields.PoolAmount,
				Params:     tt.fields.Params,
			}
			if err := gs.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GenesisState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
