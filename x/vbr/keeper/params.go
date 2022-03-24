package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the total set params
func (k Keeper) GetParamSet(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of params
func (k Keeper) SetParamSet(ctx sdk.Context, params types.Params) {
	// consider validating params, since SetParamSet could panic!

	k.paramSpace.SetParamSet(ctx, &params)
}
