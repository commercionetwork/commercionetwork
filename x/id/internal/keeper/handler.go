package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		case MsgRequestDidDeposit:
			return handleMsgRequestDidDeposit(ctx, keeper, msg)
		case MsgInvalidateDidDepositRequest:
			return handleMsgInvalidateDidDepositRequest(ctx, keeper, govKeeper, msg)
		case MsgRequestDidPowerUp:
			return handleMsgRequestDidPowerUp(ctx, keeper, msg)
		case MsgInvalidateDidPowerUpRequest:
			return handleMsgInvalidateDidPowerUpRequest(ctx, keeper, govKeeper, msg)
		case MsgMoveDeposit:
			return handleMsgWithdrawDeposit(ctx, keeper, govKeeper, msg)
		case MsgPowerUpDid:
			return handleMsgPowerUpDid(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// handleMsgSetIdentity allows to handle a MsgSetIdentity checking that the user that wants to set an identity is
// the real owner of that identity.
// If the user is not allowed to use that identity, returns an error.
func handleMsgSetIdentity(ctx sdk.Context, keeper Keeper, msg types.MsgSetIdentity) sdk.Result {
	if err := keeper.SaveDidDocument(ctx, types.DidDocument(msg)); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

func handleMsgRequestDidDeposit(ctx sdk.Context, keeper Keeper, msg MsgRequestDidDeposit) sdk.Result {

	// Set the initial status to nil and save it
	request := DidDepositRequest(msg)
	request.Status = nil

	if err := keeper.StoreDidDepositRequest(ctx, request); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgInvalidateDidDepositRequest(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper,
	msg MsgInvalidateDidDepositRequest) sdk.Result {

	// Check the status
	if msg.Status.Type != StatusRejected && msg.Status.Type != StatusCanceled {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid status: %s", msg.Status.Type)).Result()
	}

	// Check the signer if status is approved or rejected
	validGovernment := govKeeper.GetGovernmentAddress(ctx).Equals(msg.Editor)
	if msg.Status.Type == StatusRejected && !validGovernment {
		msg := fmt.Sprintf("Cannot set status of type %s without being the government", msg.Status.Type)
		return sdk.ErrInvalidAddress(msg).Result()
	}

	// Get the existing request
	existing, found := keeper.GetDidDepositRequestByProof(ctx, msg.DepositProof)
	if !found {
		msg := fmt.Sprintf("Did deposit request with proof %s not found", msg.DepositProof)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Check the signer if status is canceled
	if msg.Status.Type == StatusCanceled && !existing.FromAddress.Equals(msg.Editor) {
		return sdk.ErrInvalidAddress("Cannot edit this request without being the original poster").Result()
	}

	// Check that the existing request does not have a status set yet
	if existing.Status != nil {
		msg := fmt.Sprintf("Did deposit request with proof %s already has a valid status", existing.Proof)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Change the status, return any result
	if err := keeper.ChangeDepositRequestStatus(ctx, msg.DepositProof, msg.Status); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// ----------------------------
// --- Did power up requests
// ----------------------------

func handleMsgRequestDidPowerUp(ctx sdk.Context, keeper Keeper, msg MsgRequestDidPowerUp) sdk.Result {

	// Set the initial status to nil and save it
	request := DidPowerUpRequest(msg)
	request.Status = nil

	if err := keeper.StorePowerUpRequest(ctx, request); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgInvalidateDidPowerUpRequest(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper,
	msg MsgInvalidateDidPowerUpRequest) sdk.Result {

	// Check the status
	if msg.Status.Type != StatusRejected && msg.Status.Type != StatusCanceled {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid status: %s", msg.Status.Type)).Result()
	}

	// Check the signer if status is approved or rejected
	validGovernment := govKeeper.GetGovernmentAddress(ctx).Equals(msg.Editor)
	if msg.Status.Type == StatusRejected && !validGovernment {
		msg := fmt.Sprintf("Cannot set status of type %s without being the government", msg.Status.Type)
		return sdk.ErrInvalidAddress(msg).Result()
	}

	// Get the existing request
	existing, found := keeper.GetPowerUpRequestByProof(ctx, msg.PowerUpProof)
	if !found {
		msg := fmt.Sprintf("Did power up request with proof %s not found", msg.PowerUpProof)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Check the signer if status is canceled
	if msg.Status.Type == StatusCanceled && !existing.Claimant.Equals(msg.Editor) {
		return sdk.ErrInvalidAddress("Cannot edit this request without being the original poster").Result()
	}

	// Check that the existing request does not have a status set yet
	if existing.Status != nil {
		msg := fmt.Sprintf("Did power up request with proof %s already has a valid status", existing.Proof)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Change the status, return any result
	if err := keeper.ChangePowerUpRequestStatus(ctx, msg.PowerUpProof, msg.Status); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// ------------------------
// --- Deposits handling
// ------------------------

func handleMsgWithdrawDeposit(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg MsgMoveDeposit) sdk.Result {

	// Validate the signer
	if !govKeeper.GetGovernmentAddress(ctx).Equals(msg.Signer) {
		msg := fmt.Sprintf("Invalid signer, must be government: %s", msg.Signer)
		return sdk.ErrInvalidAddress(msg).Result()
	}

	// Get the existing request
	existing, found := keeper.GetDidDepositRequestByProof(ctx, msg.DepositProof)
	if !found {
		msg := fmt.Sprintf("Deposit request with proof %s not found", msg.DepositProof)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Check that the existing request does not have a status set yet
	if existing.Status != nil {
		msg := fmt.Sprintf("Did deposit request with proof %s already has a valid status", existing.Proof)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Move the deposit amount
	if err := keeper.DepositIntoPool(ctx, existing.FromAddress, existing.Amount); err != nil {
		return err.Result()
	}

	// Update the request
	status := RequestStatus{Type: StatusApproved}
	if err := keeper.ChangeDepositRequestStatus(ctx, existing.Proof, status); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func handleMsgPowerUpDid(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg MsgPowerUpDid) sdk.Result {

	// Validate the signer
	if !govKeeper.GetGovernmentAddress(ctx).Equals(msg.Signer) {
		msg := fmt.Sprintf("Invalid signer, must be government: %s", msg.Signer)
		return sdk.ErrInvalidAddress(msg).Result()
	}

	// Get the existing references
	references := keeper.GetHandledPowerUpRequestsReferences(ctx)
	if references.Contains(msg.ActivationReference) {
		msg := fmt.Sprintf("Power up with reference %s already handled", msg.ActivationReference)
		return sdk.ErrUnknownRequest(msg).Result()
	}

	// Move the deposit amount
	if err := keeper.FundAccount(ctx, msg.Recipient, msg.Amount); err != nil {
		return err.Result()
	}

	// Set the request as handled
	keeper.SetPowerUpRequestHandled(ctx, msg.ActivationReference)

	return sdk.Result{}
}
