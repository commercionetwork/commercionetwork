package keeper

import (
	"context"
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	eventIncrementBlockRewardsPool = "increment_block_rewards_pool"
	eventSetParams                 = "new_params"
)

func (k msgServer) IncrementBlockRewardsPool(goCtx context.Context, msg *types.MsgIncrementBlockRewardsPool) (*types.MsgIncrementBlockRewardsPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	//return &types.MsgIncrementBlockRewardsPoolResponse{}, nil
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

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventIncrementBlockRewardsPool,
		sdk.NewAttribute("funder", msg.Funder),
		sdk.NewAttribute("amount", msg.Amount.String()),
	))
	ctypes.EmitCommonEvents(ctx, msg.Funder)

	logger := k.Logger(ctx)
	logger.Debug("Block reward pool successfully increased")

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
	if msg.DistrEpochIdentifier != types.EpochDay && msg.DistrEpochIdentifier != types.EpochWeek && msg.DistrEpochIdentifier != types.EpochMinute && msg.DistrEpochIdentifier != types.EpochHour {
		return &types.MsgSetParamsResponse{}, sdkErr.Wrap(sdkErr.ErrInvalidType, fmt.Sprintf("invalid epoch identifier: %s", msg.DistrEpochIdentifier))
	}
	if msg.EarnRate.IsNegative() {
		return &types.MsgSetParamsResponse{}, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("invalid vbr earn rate: %s", msg.EarnRate))
	}
	params := types.Params{
		DistrEpochIdentifier: msg.DistrEpochIdentifier,
		EarnRate:             msg.EarnRate,
	}
	k.SetParamSet(ctx, params)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSetParams,
		sdk.NewAttribute("government", msg.Government),
		sdk.NewAttribute("distr_epoch_identifier", msg.DistrEpochIdentifier),
		sdk.NewAttribute("earn_rate", msg.EarnRate.String()),
	))
	ctypes.EmitCommonEvents(ctx, msg.Government)

	logger := k.Logger(ctx)
	logger.Debug("Params successfully set up")

	return &types.MsgSetParamsResponse{}, nil
}
