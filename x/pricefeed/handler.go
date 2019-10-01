package pricefeed

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetPrice:
			return handleMsgSetPrice(ctx, keeper, msg)
		case MsgAddOracle:
			return handleMsgAddOracle(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg MsgSetPrice) sdk.Result {
	// Check the signer
	if !keeper.IsOracle(ctx, msg.Oracle) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("%s is not an oracle", msg.Oracle.String())).Result()
	}

	// Set the raw price
	if err := keeper.SetRawPrice(ctx, RawPrice(msg)); err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}
	return sdk.Result{}
}

func handleMsgAddOracle(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg MsgAddOracle) sdk.Result {
	gov := govKeeper.GetGovernmentAddress(ctx)

	// Someone who's not the government is trying to add an oracle
	if !(gov.Equals(msg.Signer)) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("%s hasn't the rights to add an oracle", msg.Signer)).Result()
	}

	keeper.AddOracle(ctx, msg.Oracle)
	return sdk.Result{}
}
