package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgOpenCdp:
			return handleMsgOpenCdp(ctx, keeper, msg)
		case types.MsgCloseCdp:
			return handleMsgCloseCdp(ctx, keeper, msg)
		case types.MsgSetCdpCollateralRate:
			return handleMsgSetCdpCollateralRate(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgOpenCdp(ctx sdk.Context, keeper Keeper, msg types.MsgOpenCdp) (*sdk.Result, error) {
	err := keeper.OpenCdp(ctx, msg.Depositor, msg.DepositedAmount)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Log: "Cdp opened successfully"}, nil
}

func handleMsgCloseCdp(ctx sdk.Context, keeper Keeper, msg types.MsgCloseCdp) (*sdk.Result, error) {
	err := keeper.CloseCdp(ctx, msg.Signer, msg.Timestamp)
	if err != nil {
		return nil, err
	}

	return &sdk.Result{Log: "Cdp closed successfully"}, nil
}

func handleMsgSetCdpCollateralRate(ctx sdk.Context, keeper Keeper, msg types.MsgSetCdpCollateralRate) (*sdk.Result, error) {
	gov := keeper.govKeeper.GetGovernmentAddress(ctx)
	if !(gov.Equals(msg.Signer)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set collateral rate", msg.Signer))
	}
	if err := keeper.SetCollateralRate(ctx, msg.CdpCollateralRate); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	return &sdk.Result{Log: fmt.Sprintf("Cdp collateral rate changed successfully to %s", msg.CdpCollateralRate)}, nil
}
