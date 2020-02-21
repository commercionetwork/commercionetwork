package keeper

import (
	"errors"
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
)

func TestValidMsg_StoreDoc(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	handler := NewHandler(k, govK)
	msgSetIdentity := types.MsgSetIdentity(TestDidDocument)
	_, err := handler(ctx, msgSetIdentity)

	require.NoError(t, err)
}

func TestInvalidMsg(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	tm := sdk.NewTestMsg()
	handler := NewHandler(k, govK)
	_, err := handler(ctx, tm)

	require.Error(t, err)
	require.Equal(t, fmt.Sprintf("unknown request: Unrecognized %s message type: %s", types.ModuleName, tm.Type()), err.Error())
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

var msgRequestDidDeposit = types.MsgRequestDidDeposit{
	Recipient:     TestDidDepositRequest.Recipient,
	Amount:        TestDidDepositRequest.Amount,
	Proof:         TestDidDepositRequest.Proof,
	EncryptionKey: TestDidDepositRequest.EncryptionKey,
	FromAddress:   TestDidDepositRequest.FromAddress,
}

func Test_handleMsgRequestDidDeposit_NewRequest(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msgRequestDidDeposit)
	require.NoError(t, err)

	stored, err := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, TestDidDepositRequest, stored)
}

func Test_handleMsgRequestDidDeposit_ExistingRequest(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msgRequestDidDeposit)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
}

func Test_handleMsgChangeDidDepositRequestStatus_Approved_ReturnsError(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := types.RequestStatus{Type: types.StatusApproved, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.Recipient)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), msg.Status.Type)
}

func Test_handleMsgChangeDidDepositRequestStatus_Rejected_InvalidGovernment(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := types.RequestStatus{Type: types.StatusRejected, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.Recipient)
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), msg.Status.Type)
}

func Test_handleMsgChangeDidDepositRequestStatus_Rejected_ValidGovernment(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := types.RequestStatus{Type: types.StatusRejected, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)
}

func Test_handleMsgChangeDidDepositRequestStatus_Canceled_InvalidAddress(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	addr, _ := sdk.AccAddressFromBech32("cosmos1elzra8xnfqhqg2dh5ae9x45tnmud5wazkp92r9")
	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, addr)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), "poster")
}

func Test_handleMsgChangeDidDepositRequestStatus_Cancel_ValidAddress(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)

	stored, err := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, status, *stored.Status)
}

func Test_handleMsgChangeDidDepositRequestStatus_StatusAlreadySet(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	request := types.DidDepositRequest{
		Status:        &types.RequestStatus{Type: types.StatusApproved, Message: ""},
		Recipient:     TestDidDepositRequest.Recipient,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		FromAddress:   TestDidDepositRequest.FromAddress,
	}
	_ = k.StoreDidDepositRequest(ctx, request)

	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "status")
}

func Test_handleMsgChangeDidDepositRequestStatus_AllGood(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)

	stored, err := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, status, *stored.Status)
}

// ----------------------------
// --- Did power up requests
// --------------------------

var msgRequestDidPowerUp = types.MsgRequestDidPowerUp{
	Claimant:      TestDidPowerUpRequest.Claimant,
	Amount:        TestDidPowerUpRequest.Amount,
	Proof:         TestDidPowerUpRequest.Proof,
	EncryptionKey: TestDidPowerUpRequest.EncryptionKey,
}

func Test_handleMsgRequestDidPowerUp_NewRequest(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msgRequestDidPowerUp)
	require.NoError(t, err)

	stored, err := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, TestDidPowerUpRequest, stored)
}

func Test_handleMsgRequestDidPowerUp_ExistingRequest(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msgRequestDidPowerUp)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Approved_ReturnsError(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := types.RequestStatus{Type: types.StatusApproved, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), msg.Status.Type)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Rejected_InvalidGovernment(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := types.RequestStatus{Type: types.StatusRejected, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), msg.Status.Type)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Rejected_ValidGovernment(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := types.RequestStatus{Type: types.StatusRejected, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Canceled_InvalidAddress(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	addr, _ := sdk.AccAddressFromBech32("cosmos1elzra8xnfqhqg2dh5ae9x45tnmud5wazkp92r9")
	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, addr)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), "poster")
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Cancel_ValidAddress(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)

	stored, err := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, status, *stored.Status)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_StatusAlreadySet(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	request := types.DidPowerUpRequest{
		Status:        &types.RequestStatus{Type: types.StatusApproved, Message: ""},
		Amount:        TestDidPowerUpRequest.Amount,
		Proof:         TestDidPowerUpRequest.Proof,
		EncryptionKey: TestDidPowerUpRequest.EncryptionKey,
		Claimant:      TestDidPowerUpRequest.Claimant,
	}
	_ = k.StorePowerUpRequest(ctx, request)

	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "status")
}

func Test_handleMsgChangeDidPowerUpRequestStatus_AllGood(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := types.RequestStatus{Type: types.StatusCanceled, Message: ""}
	msg := types.NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)

	stored, err := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, status, *stored.Status)
}

// ------------------------
// --- Deposits handling
// ------------------------

func Test_handleMsgWithdrawDeposit_InvalidGovernment(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	msg := types.NewMsgMoveDeposit("", TestDidDepositRequest.FromAddress)
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), "government")
}

func Test_handleMsgWithdrawDeposit_InvalidRequestProof(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	msg := types.NewMsgMoveDeposit("", govK.GetGovernmentAddress(ctx))
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "not found")
}

func Test_handleMsgWithdrawDeposit_RequestAlreadyHasAStatus(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	request := types.DidDepositRequest{
		Status: &types.RequestStatus{
			Type:    "accepted",
			Message: "",
		},
		Recipient:     TestDidDepositRequest.Recipient,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		FromAddress:   TestDidDepositRequest.FromAddress,
	}
	_ = k.StoreDidDepositRequest(ctx, request)

	msg := types.NewMsgMoveDeposit(request.Proof, govK.GetGovernmentAddress(ctx))
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "already has a valid status")
}

func Test_handleMsgWithdrawDeposit_AllGood(t *testing.T) {
	_, ctx, _, bK, govK, k := SetupTestInput()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)
	_ = bK.SetCoins(ctx, TestDidDepositRequest.FromAddress, TestDidDepositRequest.Amount)

	msg := types.NewMsgMoveDeposit(TestDidDepositRequest.Proof, govK.GetGovernmentAddress(ctx))
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)
	require.NoError(t, err)

	// Check the balances
	require.Equal(t, TestDidDepositRequest.Amount, k.GetPoolAmount(ctx))
	require.Empty(t, bK.GetCoins(ctx, TestDidDepositRequest.FromAddress))
	require.Empty(t, bK.GetCoins(ctx, TestDidDepositRequest.Recipient))

	// Check the request
	request, _ := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	require.NotNil(t, request.Status)
	require.Equal(t, types.StatusApproved, request.Status.Type)
}

func Test_handleMsgPowerUpDid_InvalidGovernment(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	msg := types.MsgPowerUpDid{
		Recipient:           TestDidPowerUpRequest.Claimant,
		Amount:              TestDidPowerUpRequest.Amount,
		ActivationReference: "xxxxxx",
		Signer:              TestDidPowerUpRequest.Claimant,
	}
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrInvalidAddress))
	require.Contains(t, err.Error(), "government")
}

func Test_handleMsgPowerUpDid_ReferenceAlreadyPresent(t *testing.T) {
	_, ctx, _, _, govK, k := SetupTestInput()

	reference := "xxxxxx"
	k.SetHandledPowerUpRequestsReference(ctx, reference)

	msg := types.MsgPowerUpDid{
		Recipient:           TestDidPowerUpRequest.Claimant,
		Amount:              TestDidPowerUpRequest.Amount,
		ActivationReference: reference,
		Signer:              govK.GetGovernmentAddress(ctx),
	}
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "already handled")
}

func Test_handleMsgPowerUpDid_AllGood(t *testing.T) {
	_, ctx, _, bK, govK, k := SetupTestInput()

	msg := types.MsgPowerUpDid{
		Recipient:           TestDidPowerUpRequest.Claimant,
		Amount:              TestDidPowerUpRequest.Amount,
		ActivationReference: "test-reference",
		Signer:              govK.GetGovernmentAddress(ctx),
	}

	k.supplyKeeper.SetSupply(ctx, supply.NewSupply(msg.Amount))
	_ = bK.SetCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName), msg.Amount)
	handler := NewHandler(k, govK)
	_, err := handler(ctx, msg)

	require.NoError(t, err)

	// Check the balances
	require.Equal(t, msg.Amount, bK.GetCoins(ctx, msg.Recipient))
	require.Empty(t, k.GetPoolAmount(ctx))

	// Check the request
	require.True(t, k.GetHandledPowerUpRequestsReferences(ctx).Contains(msg.ActivationReference))
}
