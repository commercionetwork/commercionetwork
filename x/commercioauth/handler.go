package commercioauth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCreateAccount:
			return handleCreateAccount(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized commerciodocs message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// ----------------------------------
// --- CreteAccount
// ----------------------------------

func handleCreateAccount(ctx sdk.Context, keeper Keeper, msg MsgCreateAccount) sdk.Result {

	// Create the account
	err := keeper.RegisterAccount(ctx, msg.Address, msg.KeyType, msg.KeyValue)
	if err != nil {
		panic(err)
	}

	return sdk.Result{}
}
