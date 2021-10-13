package v3_0_0

import (
	v220did "github.com/commercionetwork/commercionetwork/x/did/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220did.GenesisState) *types.GenesisState {

	return &types.GenesisState{}

}
