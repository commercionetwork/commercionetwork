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
	_, ctx, govK, k := TestSetup()

	msg := NewMsgSetIdentity(TestOwnerAddress, TestDidDocument)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

// ------------------
// --- Pairwise Did
// ------------------

func Test_handleMsgRequestDidDeposit_NewRequest(t *testing.T) {
	_, ctx, govK, k := TestSetup()

	msg := NewMsgRequestDidDeposit(TestDidDepositRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, TestDidDepositRequest, stored)
}

func Test_handleMsgRequestDidDeposit_NewRequest_ExistingStatusIsReplaced(t *testing.T) {
	_, ctx, govK, k := TestSetup()

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
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	msg := NewMsgRequestDidDeposit(TestDidDepositRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleMsgChangeDidDepositRequestStatus_Approved_InvalidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusApproved, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, TestDidDepositRequest.Recipient)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidDepositRequestStatus_Approved_ValidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusApproved, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func Test_handleMsgChangeDidDepositRequestStatus_Rejected_InvalidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, TestDidDepositRequest.Recipient)
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidDepositRequestStatus_Rejected_ValidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func Test_handleMsgChangeDidDepositRequestStatus_Canceled_InvalidAddress(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	addr, _ := sdk.AccAddressFromBech32("cosmos1elzra8xnfqhqg2dh5ae9x45tnmud5wazkp92r9")
	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, addr)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, "poster")
}

func Test_handleMsgChangeDidDepositRequestStatus_Cancel_ValidAddress(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

func Test_handleMsgChangeDidDepositRequestStatus_StatusAlreadySet(t *testing.T) {
	_, ctx, govK, k := TestSetup()

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
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "status")
}

func Test_handleMsgChangeDidDepositRequestStatus_AllGood(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StoreDidDepositRequest(ctx, TestDidDepositRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidDepositRequestStatus(status, TestDidDepositRequest.Proof, TestDidDepositRequest.FromAddress)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

func Test_handleMsgRequestDidPowerUp_NewRequest(t *testing.T) {
	_, ctx, govK, k := TestSetup()

	msg := NewMsgRequestDidPowerUp(TestDidPowerUpRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)
	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, TestDidPowerUpRequest, stored)
}

func Test_handleMsgRequestDidPowerUp_NewRequest_ExistingStatusIsReplaced(t *testing.T) {
	_, ctx, govK, k := TestSetup()

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
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	msg := NewMsgRequestDidPowerUp(TestDidPowerUpRequest)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Approved_InvalidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusApproved, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Approved_ValidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusApproved, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Rejected_InvalidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)
	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, msg.Status.Type)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Rejected_ValidGovernment(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusRejected, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, govK.GetGovernmentAddress(ctx))

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Canceled_InvalidAddress(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	addr, _ := sdk.AccAddressFromBech32("cosmos1elzra8xnfqhqg2dh5ae9x45tnmud5wazkp92r9")
	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, addr)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
	assert.Contains(t, res.Log, "poster")
}

func Test_handleMsgChangeDidPowerUpRequestStatus_Cancel_ValidAddress(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}

func Test_handleMsgChangeDidPowerUpRequestStatus_StatusAlreadySet(t *testing.T) {
	_, ctx, govK, k := TestSetup()

	request := DidPowerUpRequest{
		Status:        &RequestStatus{Type: StatusApproved, Message: ""},
		Amount:        TestDidPowerUpRequest.Amount,
		Proof:         TestDidPowerUpRequest.Proof,
		EncryptionKey: TestDidPowerUpRequest.EncryptionKey,
		Claimant:      TestDidPowerUpRequest.Claimant,
	}
	_ = k.StorePowerUpRequest(ctx, request)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "status")
}

func Test_handleMsgChangeDidPowerUpRequestStatus_AllGood(t *testing.T) {
	_, ctx, govK, k := TestSetup()
	_ = k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)

	status := RequestStatus{Type: StatusCanceled, Message: ""}
	msg := NewMsgChangeDidPowerUpRequestStatus(status, TestDidPowerUpRequest.Proof, TestDidPowerUpRequest.Claimant)

	handler := NewHandler(k, govK)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, status, *stored.Status)
}
