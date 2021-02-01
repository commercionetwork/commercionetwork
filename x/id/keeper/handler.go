package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	government "github.com/commercionetwork/commercionetwork/x/government/keeper"
	"github.com/commercionetwork/commercionetwork/x/id/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		case types.MsgRequestDidPowerUp:
			return handleMsgRequestDidPowerUp(ctx, keeper, msg)
		case types.MsgChangePowerUpStatus:
			return handleMsgChangePowerUpStatus(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgSetIdentity allows to handle a MsgSetIdentity checking that the user that wants to set an identity is
// the real owner of that identity.
// If the user is not allowed to use that identity, returns an error.
func handleMsgSetIdentity(ctx sdk.Context, keeper Keeper, msg types.MsgSetIdentity) (*sdk.Result, error) {
	if err := keeper.SaveDidDocument(ctx, types.DidDocument(msg)); err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Set identity successfully added"}, nil
}

// ----------------------------
// --- Did power up requests
// ----------------------------

func handleMsgRequestDidPowerUp(ctx sdk.Context, keeper Keeper, msg types.MsgRequestDidPowerUp) (*sdk.Result, error) {

	// Crete the request
	request := types.DidPowerUpRequest{
		Claimant: msg.Claimant,
		Amount:   msg.Amount,
		Proof:    msg.Proof,
		ID:       msg.ID,
		ProofKey: msg.ProofKey,
	}

	if err := keeper.StorePowerUpRequest(ctx, request); err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Power up request successfully added"}, nil
}

// handleMsgChangePowerUpStatus marks the PowerUp request identified by the activation reference as handled successfully.
func handleMsgChangePowerUpStatus(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper, msg types.MsgChangePowerUpStatus) (*sdk.Result, error) {
	// Check the signer if status is approved or rejected
	if !govKeeper.GetTumblerAddress(ctx).Equals(msg.Signer) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Cannot set request as handled without being the tumbler")
	}

	// Get the existing request
	existing, err := keeper.GetPowerUpRequestByID(ctx, msg.PowerUpID)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	// Check the signer if status is canceled
	if !existing.Claimant.Equals(msg.Recipient) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Cannot edit this request without being the original poster")
	}

	// Check that the existing request does not have a status set yet
	if existing.Status != nil {
		msg := fmt.Sprintf("Did power up request with id %s already has a valid status", existing.Proof)
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	status := types.NewRequestStatus(msg.Status.Type, msg.Status.Message)
	// Change the status, return any result
	if err := keeper.ChangePowerUpRequestStatus(ctx, msg.PowerUpID, status); err != nil {
		return nil, err
	}

	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Power up status successfully changed"}, nil
}
