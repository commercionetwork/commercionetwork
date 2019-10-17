package id

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ------------------
// --- Identities
// ------------------

func Test_handleMsgSetIdentity(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	msg := NewMsgSetIdentity(TestOwnerAddress, TestDidDocument)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

func Test_handleMsgRequestDidDeposit_NewRequest(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	msg := NewMsgRequestDidDeposit(TestDidDepositRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, TestDidDepositRequest, stored)
}

func Test_handleMsgRequestDidDeposit_NewRequest_ExistingStatusIsReplaced(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	request := DidDepositRequest{
		Status:        &RequestStatus{Type: StatusApproved},
		Recipient:     TestDidDepositRequest.Recipient,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		FromAddress:   TestDidDepositRequest.FromAddress,
	}
	msg := NewMsgRequestDidDeposit(request)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, request.Proof)
	assert.True(t, found)
	assert.Nil(t, stored.Status)
}

func Test_handleMsgRequestDidDeposit_ExistingRequest(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	msg := NewMsgRequestDidDeposit(TestDidDepositRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleMsgChangeDidDepositRequestStatus_Approved_ReturnsError(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusApproved, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.Recipient)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidDepositRequestStatus_Rejected_InvalidGovernment(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.Recipient)
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidDepositRequestStatus_Rejected_ValidGovernment(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func Test_handleMsgChangeDidDepositRequestStatus_Canceled_InvalidAddress(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	addr, _ := sdk.AccAddressFromBech32("cosmos1elzra8xnfqhqg2dh5ae9x45tnmud5wazkp92r9")
	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, addr)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, "poster")
}

func Test_handleMsgChangeDidDepositRequestStatus_Cancel_ValidAddress(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

func Test_handleMsgChangeDidDepositRequestStatus_StatusAlreadySet(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	request := DidDepositRequest{
		Status:        &RequestStatus{Type: StatusApproved, Message: ""},
		Recipient:     TestDidDepositRequest.Recipient,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		FromAddress:   TestDidDepositRequest.FromAddress,
	}
	_ = k.StoreDidDepositRequest(ctx, request)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "status")
}

func Test_handleMsgChangeDidDepositRequestStatus_AllGood(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidDepositRequest(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

// ----------------------------
// --- Did power up requests
// --------------------------

func Test_handleMsgRequestDidPowerUp_NewRequest(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	msg := NewMsgRequestDidPowerUp(TestDidPowerUpRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, TestDidPowerUpRequest, stored)
}

func Test_handleMsgRequestDidPowerUp_NewRequest_ExistingStatusIsReplaced(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	request := DidPowerUpRequest{
		Status:        &RequestStatus{Type: StatusApproved},
		Amount:        TestDidPowerUpRequest.Amount,
		Proof:         TestDidPowerUpRequest.Proof,
		EncryptionKey: TestDidPowerUpRequest.EncryptionKey,
		Claimant:      TestDidPowerUpRequest.Claimant,
	}
	msg := NewMsgRequestDidPowerUp(request)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, request.Proof)
	assert.True(t, found)
	assert.Nil(t, stored.Status)
}

func Test_handleMsgRequestDidPowerUp_ExistingRequest(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	msg := NewMsgRequestDidPowerUp(TestDidPowerUpRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Approved_ReturnsError(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusApproved, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Rejected_InvalidGovernment(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Rejected_ValidGovernment(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Canceled_InvalidAddress(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	addr, _ := sdk.AccAddressFromBech32("cosmos1elzra8xnfqhqg2dh5ae9x45tnmud5wazkp92r9")
	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, addr)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, "poster")
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Cancel_ValidAddress(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_StatusAlreadySet(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	request := DidPowerUpRequest{
		Status:        &RequestStatus{Type: StatusApproved, Message: ""},
		Amount:        TestDidPowerUpRequest.Amount,
		Proof:         TestDidPowerUpRequest.Proof,
		EncryptionKey: TestDidPowerUpRequest.EncryptionKey,
		Claimant:      TestDidPowerUpRequest.Claimant,
	}
	_ = k.StorePowerUpRequest(ctx, request)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "status")
}

func Test_handleMsgChangeDidPowerUpRequestStatus_AllGood(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgInvalidateDidPowerUpRequest(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

// ------------------------
// --- Deposits handling
// ------------------------

func Test_handleMsgWithdrawDeposit_InvalidGovernment(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	msg := NewMsgMoveDeposit("", TestDidDepositRequest.FromAddress)
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, "government")
}

func Test_handleMsgWithdrawDeposit_InvalidRequestProof(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	msg := NewMsgMoveDeposit("", govK.GetGovernmentAddress(ctx))
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "not found")
}

func Test_handleMsgWithdrawDeposit_RequestAlreadyHasAStatus(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	request := DidDepositRequest{
		Status: &RequestStatus{
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

	msg := NewMsgMoveDeposit(request.Proof, govK.GetGovernmentAddress(ctx))
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "already has a valid status")
}

func Test_handleMsgWithdrawDeposit_AllGood(t *testing.T) {
	_, ctx, govK, bk, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)
	_ = bk.SetCoins(ctx, TestDidDepositRequest.FromAddress, TestDidDepositRequest.Amount)

	msg := NewMsgMoveDeposit(TestDidDepositRequest.Proof, govK.GetGovernmentAddress(ctx))
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	// Check the balances
	assert.Equal(t, TestDidDepositRequest.Amount, k.GetPoolAmount(ctx))
	assert.Empty(t, bk.GetCoins(ctx, TestDidDepositRequest.FromAddress))
	assert.Empty(t, bk.GetCoins(ctx, TestDidDepositRequest.Recipient))

	// Check the request
	request, _ := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.NotNil(t, request.Status)
	assert.Equal(t, StatusApproved, request.Status.Type)
}

func Test_handleMsgPowerUpDid_InvalidGovernment(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	msg := MsgPowerUpDid{
		Recipient:           TestDidPowerUpRequest.Claimant,
		Amount:              TestDidPowerUpRequest.Amount,
		ActivationReference: "xxxxxx",
		Signer:              TestDidPowerUpRequest.Claimant,
	}
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, "government")
}

func Test_handleMsgPowerUpDid_ReferenceAlreadyPresent(t *testing.T) {
	_, ctx, govK, _, k := TestSetup()

	reference := "xxxxxx"
	k.SetHandledPowerUpRequestsReferences(ctx, []string{reference})

	msg := MsgPowerUpDid{
		Recipient:           TestDidPowerUpRequest.Claimant,
		Amount:              TestDidPowerUpRequest.Amount,
		ActivationReference: reference,
		Signer:              govK.GetGovernmentAddress(ctx),
	}
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "already handled")
}

func Test_handleMsgPowerUpDid_AllGood(t *testing.T) {
	_, ctx, govK, bk, k := TestSetup()

	msg := MsgPowerUpDid{
		Recipient:           TestDidPowerUpRequest.Claimant,
		Amount:              TestDidPowerUpRequest.Amount,
		ActivationReference: "test-reference",
		Signer:              govK.GetGovernmentAddress(ctx),
	}

	_ = k.SetPoolAmount(ctx, msg.Amount)
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	// Check the balances
	assert.Equal(t, msg.Amount, bk.GetCoins(ctx, msg.Recipient))
	assert.Empty(t, k.GetPoolAmount(ctx))

	// Check the request
	assert.True(t, k.GetHandledPowerUpRequestsReferences(ctx).Contains(msg.ActivationReference))
}
