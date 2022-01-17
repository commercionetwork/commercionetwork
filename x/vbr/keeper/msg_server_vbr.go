package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)
func (k msgServer) IncrementBlockRewardsPool(goCtx context.Context, msg *types.MsgIncrementBlockRewardsPool) (*types.MsgIncrementBlockRewardsPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	funderAddr, e := sdk.AccAddressFromBech32(msg.Funder)
	if e != nil {
		return nil, e
	}
	// Subtract the coins from the account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, funderAddr, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	// Set the total rewards pool
	k.SetTotalRewardPool(ctx, k.GetTotalRewardPool(ctx).Add(sdk.NewDecCoinsFromCoins(msg.Amount...)...))

	return &types.MsgIncrementBlockRewardsPoolResponse{}, nil
}

func (k msgServer) SetParams(goCtx context.Context, msg *types.MsgSetParams) (*types.MsgSetParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	gov := k.govKeeper.GetGovernmentAddress(ctx)
	msgGovAddr, e := sdk.AccAddressFromBech32(msg.Government)
	if e != nil {
		return nil, e
	}
	if !(gov.Equals(msgGovAddr)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set params", msg.Government))
	}
	if msg.DistrEpochIdentifier != types.EpochDay && msg.DistrEpochIdentifier != types.EpochWeek && msg.DistrEpochIdentifier != types.EpochMinute{
		return &types.MsgSetParamsResponse{}, sdkErr.Wrap(sdkErr.ErrInvalidType, fmt.Sprintf("invalid epoch identifier: %s", msg.DistrEpochIdentifier))
	}
	if msg.EarnRate.IsNegative() {
		return &types.MsgSetParamsResponse{}, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("invalid vbr earn rate: %s", msg.EarnRate))
	}
	params := types.Params{
			DistrEpochIdentifier: msg.DistrEpochIdentifier,
			EarnRate: msg.EarnRate,
		}
	k.SetParamSet(ctx, params)
	
	return &types.MsgSetParamsResponse{}, nil
}