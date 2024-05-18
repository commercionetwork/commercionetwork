package epochs
/*
import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
	"github.com/commercionetwork/commercionetwork/x/epochs/keeper"
	"github.com/commercionetwork/commercionetwork/x/epochs/types"
)

// NewHandler returns a handler for epochs module messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, errors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
*/