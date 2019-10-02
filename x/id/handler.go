package id

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		case MsgRequestDidDeposit:
			return handleMsgRequestDidDeposit(ctx, keeper, msg)
		case MsgChangeDidDepositRequestStatus:
			return handleMsgChangeDidDepositRequestStatus(ctx, keeper, govKeeper, msg)
		case MsgRequestDidPowerUp:
			return handleMsgRequestDidPowerUp(ctx, keeper, msg)
		case MsgChangeDidPowerUpRequestStatus:
			return handleMsgChangeDidPowerUpRequestStatus(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetIdentity(ctx sdk.Context, keeper Keeper, msg MsgSetIdentity) sdk.Result {
	keeper.SaveIdentity(ctx, msg.Owner, msg.DidDocument)
	return sdk.Result{}
}

func handleMsgRequestDidDeposit(ctx sdk.Context, keeper Keeper, msg MsgRequestDidDeposit) sdk.Result {
	if err := keeper.StoreDidDepositRequest(ctx, DidDepositRequest(msg)); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgChangeDidDepositRequestStatus(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper,
	msg MsgChangeDidDepositRequestStatus) sdk.Result {

	// Check the status
	if err := ValidateStatus(msg.Status.Type); err != nil {
		return err.Result()
	}

	// Check the signer if status is approved or rejected
	validGovernment := govKeeper.GetGovernmentAddress(ctx).Equals(msg.Editor)
	if (msg.Status.Type == StatusApproved || msg.Status.Type == StatusRejected) && !validGovernment {
		msg := fmt.Sprintf("Cannot set status of type %s without being the government", msg.Status.Type)
		return sdk.ErrInvalidAddress(msg).Result()
	}

	// Check the signer if status is canceled
	existing, found := keeper.GetDidDepositRequestByProof(ctx, msg.DepositProof)
	if !found {
		msg := fmt.Sprintf("Did deposit request with proof %s not found", msg.DepositProof)
		return sdk.ErrUnknownRequest(msg).Result()
	}
	if !existing.FromAddress.Equals(msg.Editor) {
		return sdk.ErrInvalidAddress("Cannot edit this request without being the original poster").Result()
	}

	// Change the status, return any result
	if err := keeper.ChangeDepositRequestStatus(ctx, msg.DepositProof, msg.Status); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgRequestDidPowerUp(ctx sdk.Context, keeper Keeper, msg MsgRequestDidPowerUp) sdk.Result {
	if err := keeper.StorePowerUpRequest(ctx, DidPowerUpRequest(msg)); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgChangeDidPowerUpRequestStatus(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper,
	msg MsgChangeDidPowerUpRequestStatus) sdk.Result {

	// Check the status
	if err := ValidateStatus(msg.Status.Type); err != nil {
		return err.Result()
	}

	// Check the signer if status is approved or rejected
	validGovernment := govKeeper.GetGovernmentAddress(ctx).Equals(msg.Editor)
	if (msg.Status.Type == StatusApproved || msg.Status.Type == StatusRejected) && !validGovernment {
		msg := fmt.Sprintf("Cannot set status of type %s without being the government", msg.Status.Type)
		return sdk.ErrInvalidAddress(msg).Result()
	}

	// Check the signer if status is canceled
	existing, found := keeper.GetPowerUpRequestByProof(ctx, msg.PowerUpProof)
	if !found {
		msg := fmt.Sprintf("Did power up request with proof %s not found", msg.PowerUpProof)
		return sdk.ErrUnknownRequest(msg).Result()
	}
	if !existing.Claimant.Equals(msg.Editor) {
		return sdk.ErrInvalidAddress("Cannot edit this request without being the original poster").Result()
	}

	// Change the status, return any result
	if err := keeper.ChangePowerUpRequestStatus(ctx, msg.PowerUpProof, msg.Status); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
