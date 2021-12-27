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

	msgGovernment, _ := sdk.AccAddressFromBech32(msg.Government)
	if !k.GovKeeper.GetGovernmentAddress(ctx).Equals(msgGovernment) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}

	msgTsp, _ := sdk.AccAddressFromBech32(msg.Tsp)
	membership, err := k.GetMembership(ctx, msgTsp)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Tsp %s has no membership", msg.Tsp))
	}

	if membership.MembershipType != types.MembershipTypeBlack {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Membership of Tsp %s is %s but must be %s", msg.Tsp, membership.MembershipType, types.MembershipTypeBlack))
	}

	k.AddTrustedServiceProvider(ctx, msgTsp)

	//TODO emits events
	//ctypes.EmitCommonEvents(ctx, msg.Government)

	return &types.MsgAddTspResponse{
		Tsp: msg.Tsp,
	}, nil

}

func (k msgServer) RemoveTsp(goCtx context.Context, msg *types.MsgRemoveTsp) (*types.MsgRemoveTspResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	msgGovernment, _ := sdk.AccAddressFromBech32(msg.Government)
	if !k.GovKeeper.GetGovernmentAddress(ctx).Equals(msgGovernment) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}

	msgTsp, _ := sdk.AccAddressFromBech32(msg.Tsp)
	k.RemoveTrustedServiceProvider(ctx, msgTsp)
	return &types.MsgRemoveTspResponse{
		Tsp: msg.Tsp,
	}, nil
}

func (k msgServer) DepositIntoLiquidityPool(goCtx context.Context, msg *types.MsgDepositIntoLiquidityPool) (*types.MsgDepositIntoLiquidityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	msgDepositor, _ := sdk.AccAddressFromBech32(msg.Depositor)
	if err := k.DepositIntoPool(ctx, msgDepositor, msg.Amount); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}
	return &types.MsgDepositIntoLiquidityPoolResponse{
		AmountPool: msg.Amount, // TODO response with total pool amount
	}, nil

}
