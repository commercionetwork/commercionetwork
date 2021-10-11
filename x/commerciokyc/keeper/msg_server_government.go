package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// AddTrustedServiceProvider allows to add the given signer as a trusted entity
// that can sign transactions setting an accrediter for a user.
func (k msgServer) AddTsp(goCtx context.Context, msg *types.MsgAddTsp) (*types.MsgAddTspResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.govKeeper.GetGovernmentAddress(ctx).Equals(sdk.AccAddress(msg.Government)) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}

	membership, err := k.GetMembership(ctx, sdk.AccAddress(msg.Tsp))
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Tsp %s has no membership", msg.Tsp))
	}

	if membership.MembershipType != types.MembershipTypeBlack {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Membership of Tsp %s is %s but must be %s", msg.Tsp, membership.MembershipType, types.MembershipTypeBlack))
	}

	k.AddTrustedServiceProvider(ctx, sdk.AccAddress(msg.Tsp))

	//TODO emits events
	//ctypes.EmitCommonEvents(ctx, msg.Government)

	return &types.MsgAddTspResponse{
		Tsp: msg.Tsp,
	}, nil

}

func (k msgServer) RemoveTsp(goCtx context.Context, msg *types.MsgRemoveTsp) (*types.MsgRemoveTspResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !k.govKeeper.GetGovernmentAddress(ctx).Equals(sdk.AccAddress(msg.Government)) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}

	k.RemoveTrustedServiceProvider(ctx, sdk.AccAddress(msg.Tsp))
	return &types.MsgRemoveTspResponse{
		Tsp: msg.Tsp,
	}, nil
}

func (k msgServer) DepositIntoLiquidityPool(goCtx context.Context, msg *types.MsgDepositIntoLiquidityPool) (*types.MsgDepositIntoLiquidityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := k.DepositIntoPool(ctx, sdk.AccAddress(msg.Depositor), msg.Amount); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}
	return &types.MsgDepositIntoLiquidityPoolResponse{
		msg.Amount, // TODO response with total pool amount
	}, nil

}
