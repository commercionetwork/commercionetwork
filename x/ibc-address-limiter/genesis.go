package ibc_address_limit

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/ibc-address-limiter/types"
)

// InitGenesis initializes the x/ibc-address-limiter module's state from a provided genesis
// state, which includes the parameter for the contract address.
func (i *ICS4Wrapper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	i.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the x/ibc-address-limiter module's exported genesis.
func (i *ICS4Wrapper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: types.Params{ ContractAddress: i.GetParams(ctx)},
	}
}
