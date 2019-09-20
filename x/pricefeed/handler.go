package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetPrice:
			return handleMsgSetPrice(ctx, keeper, msg)
		case MsgAddOracle:
			return handleMsgAddOracle(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetPrice(ctx sdk.Context, keeper Keeper, msg MsgSetPrice) sdk.Result {
	err := keeper.SetRawPrice(ctx, msg.Price)
	if err != nil {
		return err.Result()
	}
	return sdk.Result{}
}

func handleMsgAddOracle(ctx sdk.Context, keeper Keeper, msg MsgAddOracle) sdk.Result {

	gov := keeper.GovernmentKeeper.GetGovernmentAddress(ctx)

	//Someone who's not the government trying to add an oracle
	if !(gov.Equals(msg.Signer)) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("%s haven't got the rights to add an oracle", msg.Signer)).Result()
	}

	keeper.AddOracle(ctx, msg.Oracle)
	return sdk.Result{}
}
