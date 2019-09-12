package accreditations

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetAccrediter:
			return handleSetAccrediter(ctx, keeper, msg)
		case MsgDistributeReward:
			return handleDistributeReward(ctx, keeper, msg)
		case MsgDepositIntoLiquidityPool:
			return handleDepositIntoPool(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleSetAccrediter(ctx sdk.Context, keeper Keeper, msg MsgSetAccrediter) sdk.Result {

	// Check the signer
	if !keeper.IsTrustedSigner(ctx, msg.Signer) {
		errMsg := fmt.Sprintf("The signer %s is not trusted", msg.Signer.String())
		return sdk.ErrInvalidAddress(errMsg).Result()
	}

	// Check the accreditation
	if accreditation := keeper.GetAccreditation(ctx, msg.User); accreditation.Accrediter != nil {
		errMsg := fmt.Sprintf(
			"User %s already has an accreditation (%s)",
			msg.User.String(),
			accreditation.Accrediter.String(),
		)
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	// If everything passes the checks, set the accreditation
	if err := keeper.SetAccrediter(ctx, msg.User, msg.Accrediter); err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}

	return sdk.Result{}
}

func handleDistributeReward(ctx sdk.Context, keeper Keeper, msg MsgDistributeReward) sdk.Result {

	// Check the accreditation
	accreditation := keeper.GetAccreditation(ctx, msg.User)
	if accreditation.Rewarded {
		return sdk.ErrUnknownRequest("Accrediter has already been rewarded").Result()
	}

	if accreditation.Accrediter == nil || !accreditation.Accrediter.Equals(msg.Accrediter) {
		errMsg := fmt.Sprintf("Accrediter of %s does not match with the given one", msg.User.String())
		return sdk.ErrInvalidAddress(errMsg).Result()
	}

	// Distribute the reward
	if err := keeper.DistributeReward(ctx, msg.Accrediter, msg.Reward, msg.User); err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}

	return sdk.Result{}
}

func handleDepositIntoPool(ctx sdk.Context, keeper Keeper, msg MsgDepositIntoLiquidityPool) sdk.Result {
	if err := keeper.DepositIntoPool(ctx, msg.Depositor, msg.Amount); err != nil {
		return sdk.ErrUnknownRequest(err.Error()).Result()
	}

	return sdk.Result{}
}
