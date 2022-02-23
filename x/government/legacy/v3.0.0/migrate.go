package v3_0_0

import (
	v220gov "github.com/commercionetwork/commercionetwork/x/government/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/government/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220gov.GenesisState) *types.GenesisState {

	govAddr := oldGenState.GovernmentAddress.String()

	return &types.GenesisState{
		GovernmentAddress: govAddr,
	}

}
