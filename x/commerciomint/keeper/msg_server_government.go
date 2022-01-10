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
	signerAccAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !(gov.Equals(signerAccAddr)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("could not set conversion rate since %s is not the government", msg.Signer))
	}
	if err := k.UpdateConversionRate(ctx, msg.Rate); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSetConversionRate,
		sdk.NewAttribute("rate", msg.Rate.String()),
	))

	return &types.MsgSetCCCConversionRateResponse{Rate: msg.Rate}, nil
}

// TODO IMPLEMENTATION
func (k msgServer) SetFreezePeriod(goCtx context.Context, msg *types.MsgSetCCCFreezePeriod) (*types.MsgSetCCCFreezePeriodResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	gov := k.govKeeper.GetGovernmentAddress(ctx)
	// TODO MOVE TO VALIDATION
	signerAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !(gov.Equals(signerAddr)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("could not set freeze period since %s is not the government", msg.Signer))
	}
	// TODO MOVE TO VALIDATION
	freezePeriod, err := time.ParseDuration(msg.FreezePeriod)
	if err != nil {
		return nil, err
	}
	if err := k.UpdateFreezePeriod(ctx, freezePeriod); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSetFreezePeriod,
		sdk.NewAttribute("freeze_period", msg.FreezePeriod),
	))

	return &types.MsgSetCCCFreezePeriodResponse{FreezePeriod: msg.FreezePeriod}, nil
}
