package keeper

import (
	"context"
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
)

// refactor error variables names

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

// SetParams should use the Validate for Params
// (DONE) unifrom returning nil in case of error

func (k msgServer) SetParams(goCtx context.Context, msg *types.MsgSetParams) (*types.MsgSetParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	gov := k.govKeeper.GetGovernmentAddress(ctx)
	msgGovAddr, e := sdk.AccAddressFromBech32(msg.Government)
	if e != nil {
		return nil, e
	}
	if !(gov.Equals(msgGovAddr)) {
		return nil, errors.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set params", msg.Government))
	}

	params := types.NewParams(msg.DistrEpochIdentifier, msg.EarnRate)

	if err := k.SetParamSet(ctx, params); err != nil {
		return nil, errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("invalid params: %s", msg.EarnRate))
	}

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
