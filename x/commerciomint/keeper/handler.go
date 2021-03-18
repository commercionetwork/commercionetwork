package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgMintCCC:
			return handleMsgMintCCC(ctx, keeper, msg)
		case types.MsgBurnCCC:
			return handleMsgBurnCCC(ctx, keeper, msg)
		case types.MsgSetCCCConversionRate:
			return handleMsgSetCCCConversionRate(ctx, keeper, msg)
		case types.MsgSetCCCFreezePeriod:
			return handleMsgSetCCCFreezePeriod(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgMintCCC(ctx sdk.Context, keeper Keeper, msg types.MsgMintCCC) (*sdk.Result, error) {
	err := keeper.NewPosition(ctx, msg.Owner, msg.Credits, msg.ID)
	if err != nil {
		return nil, sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "cannot mint ccc, %s", err.Error())
	}
	ctypes.EmitCommonEvents(ctx, msg.Owner)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "mint successful"}, nil
}

func handleMsgBurnCCC(ctx sdk.Context, keeper Keeper, msg types.MsgBurnCCC) (*sdk.Result, error) {
	err := keeper.BurnCCC(ctx, msg.Signer, msg.ID, msg.Amount)
	if err != nil {
		return nil, err
	}
	ctypes.EmitCommonEvents(ctx, msg.Signer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "burn successful"}, nil
}

func handleMsgSetCCCConversionRate(ctx sdk.Context, keeper Keeper, msg types.MsgSetCCCConversionRate) (*sdk.Result, error) {
	gov := keeper.govKeeper.GetGovernmentAddress(ctx)
	if !(gov.Equals(msg.Signer)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set conversion rate", msg.Signer))
	}
	if err := keeper.SetConversionRate(ctx, msg.Rate); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	ctypes.EmitCommonEvents(ctx, msg.Signer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: fmt.Sprintf("conversion rate changed successfully to %s", msg.Rate)}, nil
}

func handleMsgSetCCCFreezePeriod(ctx sdk.Context, keeper Keeper, msg types.MsgSetCCCFreezePeriod) (*sdk.Result, error) {
	gov := keeper.govKeeper.GetGovernmentAddress(ctx)
	if !(gov.Equals(msg.Signer)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set conversion rate", msg.Signer))
	}
	if err := keeper.SetFreezePeriod(ctx, msg.FreezePeriod); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	ctypes.EmitCommonEvents(ctx, msg.Signer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: fmt.Sprintf("conversion rate changed successfully to %s", msg.FreezePeriod)}, nil
}
