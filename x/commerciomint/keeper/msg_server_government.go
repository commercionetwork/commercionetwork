package keeper

import (
	"context"
	"fmt"
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) SetConversionRate(goCtx context.Context, msg *types.MsgSetCCCConversionRate) (*types.MsgSetCCCConversionRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	gov := k.govKeeper.GetGovernmentAddress(ctx)
	if !(gov.Equals(sdk.AccAddress(msg.Signer))) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set conversion rate", msg.Signer))
	}
	if err := k.UpdateConversionRate(ctx, msg.Rate); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	// TODO EMITS EVENTS CORRECTLY
	/*ctypes.EmitCommonEvents(ctx, msg.Signer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: fmt.Sprintf("conversion rate changed successfully to %s", msg.Rate)}, nil
	*/

	return &types.MsgSetCCCConversionRateResponse{Rate: msg.Rate}, nil
}

// TODO IMPLEMENTATION
func (k msgServer) SetFreezePeriod(goCtx context.Context, msg *types.MsgSetCCCFreezePeriod) (*types.MsgSetCCCFreezePeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	gov := k.govKeeper.GetGovernmentAddress(ctx)
	// TODO MOVE TO VALIDATION
	if !(gov.Equals(sdk.AccAddress(msg.Signer))) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set conversion rate", msg.Signer))
	}
	// TODO MOVE TO VALIDATION
	freezDuration, err := time.ParseDuration(msg.FreezePeriod)
	if err != nil {
		return &types.MsgSetCCCFreezePeriodResponse{}, err
	}
	if err := k.UpdateFreezePeriod(ctx, freezDuration); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	// TODO EMITS EVENTS CORRECTLY

	return &types.MsgSetCCCFreezePeriodResponse{}, nil
}
