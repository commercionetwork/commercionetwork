package v1_3_4

import (
	v133tbr "github.com/commercionetwork/commercionetwork/x/vbr/legacy/v1.3.3"
)

// Migrate migrates exported state from v1.3.3 to a v1.3.4 genesis state
func Migrate(oldGenState v133tbr.GenesisState) GenesisState {
	return GenesisState{
		PoolAmount:       oldGenState.PoolAmount,
		YearlyPoolAmount: oldGenState.YearlyPoolAmount,
		YearNumber:       oldGenState.YearNumber,
	}
}
