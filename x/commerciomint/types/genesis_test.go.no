package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGenesisState_Validate(t *testing.T) {
	type fields struct {
		Positions  []*Position
		PoolAmount sdk.Coins
		Params     Params
	}
	tests := []struct {
		name         string
		genesisState func() GenesisState
		wantErr      bool
	}{
		{
			name: "valid default genesis",
			genesisState: func() GenesisState {
				return *DefaultGenesis()
			},
			wantErr: false,
		},
		{
			name: "invalid positions",
			genesisState: func() GenesisState {
				gs := *DefaultGenesis()
				invalidPosition := NewPosition(ownerAddr, sdk.NewInt(testEtp.Collateral), *testEtp.Credits, "abcd", testCreatedAt, testEtp.ExchangeRate)
				gs.Positions = []*Position{&invalidPosition}
				return gs
			},
			wantErr: true,
		},
		{
			name: "invalid params",
			genesisState: func() GenesisState {
				gs := *DefaultGenesis()
				gs.Params = NewParams(invalidConversionRate, invalidFreezePeriod)
				return gs
			},
			wantErr: true,
		},
		{
			name: "invalid pool amount",
			genesisState: func() GenesisState {
				gs := *DefaultGenesis()
				gs.PoolAmount = sdk.Coins{sdk.Coin{Denom: BondDenom, Amount: sdk.NewInt(-1)}}
				return gs
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.genesisState().Validate(); (err != nil) != tt.wantErr {
				t.Errorf("GenesisState.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
