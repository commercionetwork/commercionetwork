package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ------------------
// --- Identities
// ------------------

func TestKeeper_CreateIdentity(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)

	k.SaveIdentity(ctx, TestOwnerAddress, TestDidDocument)

	var stored types.DidDocument
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, TestDidDocument, stored)
}

func TestKeeper_EditIdentity(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	updatedDidDocument := types.DidDocument{Uri: "ddo-reference-update", ContentHash: TestDidDocument.ContentHash}
	k.SaveIdentity(ctx, TestOwnerAddress, updatedDidDocument)

	var stored types.DidDocument
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, updatedDidDocument, stored)
}

func TestKeeper_GetDidDocumentByOwner_ExistingDidDocument(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, found := k.GetDidDocumentByOwner(ctx, TestOwnerAddress)

	assert.True(t, found)
	assert.Equal(t, TestDidDocument, actual)
}

func TestKeeper_GetIdentities(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	first, _ := sdk.AccAddressFromBech32("cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second, _ := sdk.AccAddressFromBech32("cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third, _ := sdk.AccAddressFromBech32("cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth, _ := sdk.AccAddressFromBech32("cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth, _ := sdk.AccAddressFromBech32("cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

	store.Set(k.getIdentityStoreKey(first), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(second), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(third), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(fourth), cdc.MustMarshalBinaryBare(TestDidDocument))
	store.Set(k.getIdentityStoreKey(fifth), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, err := k.GetIdentities(ctx)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(actual))
	assert.Contains(t, actual, types.Identity{Owner: first, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: second, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: third, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: fourth, DidDocument: TestDidDocument})
	assert.Contains(t, actual, types.Identity{Owner: fifth, DidDocument: TestDidDocument})
}

func TestKeeper_SetIdentities(t *testing.T) {
	_, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	first, _ := sdk.AccAddressFromBech32("cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second, _ := sdk.AccAddressFromBech32("cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third, _ := sdk.AccAddressFromBech32("cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth, _ := sdk.AccAddressFromBech32("cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth, _ := sdk.AccAddressFromBech32("cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

	identities := []types.Identity{
		{Owner: first, DidDocument: TestDidDocument},
		{Owner: second, DidDocument: TestDidDocument},
		{Owner: third, DidDocument: TestDidDocument},
		{Owner: fourth, DidDocument: TestDidDocument},
		{Owner: fifth, DidDocument: TestDidDocument},
	}
	k.SetIdentities(ctx, identities)

	assert.True(t, store.Has(k.getIdentityStoreKey(first)))
	assert.True(t, store.Has(k.getIdentityStoreKey(second)))
	assert.True(t, store.Has(k.getIdentityStoreKey(third)))
	assert.True(t, store.Has(k.getIdentityStoreKey(fourth)))
	assert.True(t, store.Has(k.getIdentityStoreKey(fifth)))
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

var requestSender, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var requestRecipient, _ = sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")

func TestKeeper_StoreDidDepositRequest_NewRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	request := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
	}

	err := k.StoreDidDepositRequest(ctx, request)
	assert.Nil(t, err)

	var stored types.DidDepositRequest
	store := ctx.KVStore(k.StoreKey)
	storedBz := store.Get(k.getDepositRequestStoreKey(request.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, request, stored)
}

func TestKeeper_StoreDidDepositRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	request := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
	}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getDepositRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	err := k.StoreDidDepositRequest(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "same proof")
}

func TestKeeper_GetDidDepositRequestByProof_NonExistingRequest(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	_, found := k.GetDidDepositRequestByProof(ctx, "")
	assert.False(t, found)
}

func TestKeeper_GetDidDepositRequestByProof_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	request := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
	}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getDepositRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	stored, found := k.GetDidDepositRequestByProof(ctx, request.Proof)
	assert.True(t, found)
	assert.Equal(t, request, stored)
}

func TestKeeper_ChangeDepositRequestStatus_NonExistingRequest(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	status := types.DidDepositRequestStatus{
		Type:    "status-type",
		Message: "status-message",
	}

	err := k.ChangeDepositRequestStatus(ctx, "", status)
	assert.Error(t, err)

	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "proof")
}

func TestKeeper_ChangeDepositRequestStatus_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	request := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
	}
	store.Set(k.getDepositRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	status := types.DidDepositRequestStatus{Type: "status-type", Message: "status-message"}
	err := k.ChangeDepositRequestStatus(ctx, request.Proof, status)
	assert.Nil(t, err)

	expected := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
		Status:        &status,
	}

	var stored types.DidDepositRequest
	storedBz := store.Get(k.getDepositRequestStoreKey(request.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)
	assert.Equal(t, expected, stored)
}

func TestKeeper_GetDepositRequests_EmptyList(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	requests := k.GetDepositRequests(ctx)
	assert.Empty(t, requests)
}

func TestKeeper_GetDepositRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	request := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
	}
	store.Set(k.getDepositRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	requests := k.GetDepositRequests(ctx)
	assert.Equal(t, 1, len(requests))
	assert.Contains(t, requests, request)
}

// ----------------------------
// --- Did PowerUp requests
// ----------------------------

func TestKeeper_StorePowerUpRequest_NewRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	request := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
	}

	err := k.StorePowerUpRequest(ctx, request)
	assert.Nil(t, err)

	var stored types.DidPowerUpRequest
	store := ctx.KVStore(k.StoreKey)
	storedBz := store.Get(k.getDidPowerUpRequestStoreKey(request.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, request, stored)
}

func TestKeeper_StorePowerUpRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	request := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
	}
	store.Set(k.getDidPowerUpRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	err := k.StorePowerUpRequest(ctx, request)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "same proof")
}

func TestKeeper_GetPowerUpRequestByProof_NonExistingRequest(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	_, found := k.GetPowerUpRequestByProof(ctx, "")
	assert.False(t, found)
}

func TestKeeper_GetPowerUpRequestByProof_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	request := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
	}
	store.Set(k.getDidPowerUpRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	stored, found := k.GetPowerUpRequestByProof(ctx, request.Proof)
	assert.True(t, found)
	assert.Equal(t, request, stored)
}

func TestKeeper_ChangePowerUpRequestStatus_NonExistingRequest(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	status := types.DidPowerUpRequestStatus{
		Type:    "status-type",
		Message: "status-messsge",
	}

	err := k.ChangePowerUpRequestStatus(ctx, "", status)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "proof")
}

func TestKeeper_ChangePowerUpRequestStatus_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	request := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
	}
	store.Set(k.getDidPowerUpRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	status := types.DidPowerUpRequestStatus{
		Type:    "status-type",
		Message: "status-messsge",
	}

	err := k.ChangePowerUpRequestStatus(ctx, request.Proof, status)
	assert.Nil(t, err)

	expected := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		Status:        &status,
	}

	var stored types.DidPowerUpRequest
	storedBz := store.Get(k.getDidPowerUpRequestStoreKey(request.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)
	assert.Equal(t, expected, stored)
}

func TestKeeper_GetPowerUpRequests_EmptyList(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	requests := k.GetPowerUpRequests(ctx)
	assert.Empty(t, requests)
}

func TestKeeper_GetPowerUpRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	request := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
	}
	store.Set(k.getDidPowerUpRequestStoreKey(request.Proof), cdc.MustMarshalBinaryBare(&request))

	requests := k.GetPowerUpRequests(ctx)
	assert.Equal(t, 1, len(requests))
	assert.Contains(t, requests, request)
}
