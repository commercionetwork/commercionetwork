package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetPrice:
			return handleMsgSetPrice(ctx, keeper, msg)
		case types.MsgAddOracle:
			return handleMsgAddOracle(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (errMsg))
		}
	}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg types.MsgSetPrice) sdk.Result {
	// Check the signer
	if !keeper.IsOracle(ctx, msg.Oracle) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s is not an oracle", msg.Oracle.String()))
	}

	// Set the raw price
	if err := keeper.AddRawPrice(ctx, msg.Oracle, msg.Price); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (err.Error()))
	}
	return sdk.Result{}
}

func handleMsgAddOracle(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg types.MsgAddOracle) sdk.Result {
	gov := govKeeper.GetGovernmentAddress(ctx)

	// Someone who's not the government is trying to add an oracle
	if !(gov.Equals(msg.Signer)) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s hasn't the rights to add an oracle", msg.Signer))
	}

	keeper.AddOracle(ctx, msg.Oracle)
	return sdk.Result{}
}
