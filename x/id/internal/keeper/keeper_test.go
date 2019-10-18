package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

// ------------------
// --- Identities
// ------------------

func TestKeeper_CreateIdentity(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	err := k.SaveDidDocument(ctx, TestDidDocument)
	assert.NoError(t, err)

	var stored types.DidDocument
	store := ctx.KVStore(k.storeKey)
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, TestDidDocument, stored)
}

func TestKeeper_EditIdentity(t *testing.T) {
	cdc, ctx, aK, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	updatedDocument := types.DidDocument{}
	err := copier.Copy(&updatedDocument, &TestDidDocument)
	assert.NoError(t, err)

	account := aK.GetAccount(ctx, TestOwnerAddress)
	updatedDocument.PubKeys = types.PubKeys{
		types.PubKey{
			Id:           "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1",
			Type:         "Secp256k1VerificationKey2018",
			Controller:   TestOwnerAddress,
			PublicKeyHex: hex.EncodeToString(account.GetPubKey().Bytes()),
		},
	}

	err = k.SaveDidDocument(ctx, updatedDocument)
	assert.NoError(t, err)

	var stored types.DidDocument
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, updatedDocument, stored)
	assert.Len(t, stored.PubKeys, 1)
}

func TestKeeper_GetDidDocumentByOwner_ExistingDidDocument(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, found := k.GetDidDocumentByOwner(ctx, TestOwnerAddress)

	assert.True(t, found)
	assert.Equal(t, TestDidDocument, actual)
}

func TestKeeper_GetDidDocuments(t *testing.T) {
	cdc, ctx, aK, _, _, k := SetupTestInput()
	store := ctx.KVStore(k.storeKey)

	first := setupDidDocument(ctx, aK, "cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second := setupDidDocument(ctx, aK, "cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third := setupDidDocument(ctx, aK, "cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth := setupDidDocument(ctx, aK, "cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth := setupDidDocument(ctx, aK, "cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

	store.Set(k.getIdentityStoreKey(first.Id), cdc.MustMarshalBinaryBare(first))
	store.Set(k.getIdentityStoreKey(second.Id), cdc.MustMarshalBinaryBare(second))
	store.Set(k.getIdentityStoreKey(third.Id), cdc.MustMarshalBinaryBare(third))
	store.Set(k.getIdentityStoreKey(fourth.Id), cdc.MustMarshalBinaryBare(fourth))
	store.Set(k.getIdentityStoreKey(fifth.Id), cdc.MustMarshalBinaryBare(fifth))

	actual, err := k.GetDidDocuments(ctx)

	assert.Nil(t, err)
	assert.Equal(t, 5, len(actual))
	assert.Contains(t, actual, first)
	assert.Contains(t, actual, second)
	assert.Contains(t, actual, third)
	assert.Contains(t, actual, fourth)
	assert.Contains(t, actual, fifth)
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

func TestKeeper_StoreDidDepositRequest_NewRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	err := k.StoreDidDepositRequest(ctx, TestDidDepositRequest)
	assert.Nil(t, err)

	var stored types.DidDepositRequest
	store := ctx.KVStore(k.storeKey)
	storedBz := store.Get(k.getDepositRequestStoreKey(TestDidDepositRequest.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, TestDidDepositRequest, stored)
}

func TestKeeper_StoreDidDepositRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDepositRequestStoreKey(TestDidDepositRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidDepositRequest))

	err := k.StoreDidDepositRequest(ctx, TestDidDepositRequest)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "same proof")
}

func TestKeeper_GetDidDepositRequestByProof_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	_, found := k.GetDidDepositRequestByProof(ctx, "")
	assert.False(t, found)
}

func TestKeeper_GetDidDepositRequestByProof_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDepositRequestStoreKey(TestDidDepositRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidDepositRequest))

	stored, found := k.GetDidDepositRequestByProof(ctx, TestDidDepositRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, TestDidDepositRequest, stored)
}

func TestKeeper_ChangeDepositRequestStatus_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	status := types.RequestStatus{
		Type:    "status-type",
		Message: "status-message",
	}

	err := k.ChangeDepositRequestStatus(ctx, "", status)
	assert.Error(t, err)

	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "proof")
}

func TestKeeper_ChangeDepositRequestStatus_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDepositRequestStoreKey(TestDidDepositRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidDepositRequest))

	status := types.RequestStatus{Type: "status-type", Message: "status-message"}
	err := k.ChangeDepositRequestStatus(ctx, TestDidDepositRequest.Proof, status)
	assert.Nil(t, err)

	expected := types.DidDepositRequest{
		Recipient:     TestDidDepositRequest.Recipient,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		FromAddress:   TestDidDepositRequest.FromAddress,
		Status:        &status,
	}

	var stored types.DidDepositRequest
	storedBz := store.Get(k.getDepositRequestStoreKey(TestDidDepositRequest.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)
	assert.Equal(t, expected, stored)
}

func TestKeeper_GetDepositRequests_EmptyList(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	didDepositRequests := k.GetDepositRequests(ctx)
	assert.Empty(t, didDepositRequests)
}

func TestKeeper_GetDepositRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDepositRequestStoreKey(TestDidDepositRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidDepositRequest))

	didDepositRequests := k.GetDepositRequests(ctx)
	assert.Equal(t, 1, len(didDepositRequests))
	assert.Contains(t, didDepositRequests, TestDidDepositRequest)
}

// ----------------------------
// --- Did power up requests
// ----------------------------

func TestKeeper_StorePowerUpRequest_NewRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	err := k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)
	assert.Nil(t, err)

	var stored types.DidPowerUpRequest
	store := ctx.KVStore(k.storeKey)
	storedBz := store.Get(k.getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, TestDidPowerUpRequest, stored)
}

func TestKeeper_StorePowerUpRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	err := k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "same proof")
}

func TestKeeper_GetPowerUpRequestByProof_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	_, found := k.GetPowerUpRequestByProof(ctx, "")
	assert.False(t, found)
}

func TestKeeper_GetPowerUpRequestByProof_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	stored, found := k.GetPowerUpRequestByProof(ctx, TestDidPowerUpRequest.Proof)
	assert.True(t, found)
	assert.Equal(t, TestDidPowerUpRequest, stored)
}

func TestKeeper_ChangePowerUpRequestStatus_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	status := types.RequestStatus{
		Type:    "status-type",
		Message: "status-messsge",
	}

	err := k.ChangePowerUpRequestStatus(ctx, "", status)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "proof")
}

func TestKeeper_ChangePowerUpRequestStatus_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	status := types.RequestStatus{
		Type:    "status-type",
		Message: "status-messsge",
	}

	err := k.ChangePowerUpRequestStatus(ctx, TestDidPowerUpRequest.Proof, status)
	assert.Nil(t, err)

	expected := types.DidPowerUpRequest{
		Claimant:      TestDidPowerUpRequest.Claimant,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		Status:        &status,
	}

	var stored types.DidPowerUpRequest
	storedBz := store.Get(k.getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)
	assert.Equal(t, expected, stored)
}

func TestKeeper_GetPowerUpRequests_EmptyList(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	didPowerUpRequests := k.GetPowerUpRequests(ctx)
	assert.Empty(t, didPowerUpRequests)
}

func TestKeeper_GetPowerUpRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	didPowerUpRequests := k.GetPowerUpRequests(ctx)
	assert.Equal(t, 1, len(didPowerUpRequests))
	assert.Contains(t, didPowerUpRequests, TestDidPowerUpRequest)
}

// ------------------------
// --- Deposits handling
// ------------------------

func TestKeeper_DepositIntoPool_InvalidAmount(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	err := k.DepositIntoPool(ctx, TestDepositor, sdk.NewCoins())
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

func TestKeeper_DepositIntoPool_InsufficientFunds(t *testing.T) {
	_, ctx, _, bK, _, k := SetupTestInput()
	_ = bK.SetCoins(ctx, TestDepositor, sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)))

	err := k.DepositIntoPool(ctx, TestDepositor, sdk.NewCoins(sdk.NewInt64Coin("uatom", 1000)))
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInsufficientCoins, err.Code())
}

func TestKeeper_DepositIntoPool_ValidRequest(t *testing.T) {
	_, ctx, _, bK, _, k := SetupTestInput()
	_ = bK.SetCoins(ctx, TestDepositor, sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)))

	err := k.DepositIntoPool(ctx, TestDepositor, sdk.NewCoins(sdk.NewInt64Coin("uatom", 25)))
	assert.Nil(t, err)

	pool := k.GetPoolAmount(ctx)
	assert.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uatom", 25)), pool)
	assert.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uatom", 75)), bK.GetCoins(ctx, TestDepositor))
}

func TestKeeper_FundAccount_InvalidAmount(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	err := k.FundAccount(ctx, TestDepositor, sdk.NewCoins())
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInvalidCoins, err.Code())
}

func TestKeeper_FundAccount_InsufficientPoolFunds(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()
	_ = k.SetPoolAmount(ctx, sdk.NewCoins(sdk.NewInt64Coin("uatom", 10)))

	err := k.FundAccount(ctx, TestDepositor, sdk.NewCoins(sdk.NewInt64Coin("uatom", 1000)))
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInsufficientFunds, err.Code())
}

func TestKeeper_FundAccount_ValidRequest(t *testing.T) {
	_, ctx, _, bK, _, k := SetupTestInput()
	_ = k.SetPoolAmount(ctx, sdk.NewCoins(sdk.NewInt64Coin("uatom", 1000)))

	err := k.FundAccount(ctx, TestDepositor, sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)))
	assert.Nil(t, err)

	pool := k.GetPoolAmount(ctx)
	assert.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uatom", 900)), pool)
	assert.Equal(t, sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)), bK.GetCoins(ctx, TestDepositor))
}

func TestKeeper_SetPoolAmount_EmptyCoins(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	err := k.SetPoolAmount(ctx, nil)
	assert.Nil(t, err)
}

func TestKeeper_SetPoolAmount_NonEmptyCoins(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	pool := sdk.NewCoins(sdk.NewInt64Coin("uatom", 100))
	err := k.SetPoolAmount(ctx, pool)
	assert.Nil(t, err)

	var stored sdk.Coins
	store := ctx.KVStore(k.storeKey)
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.DepositsPoolStoreKey)), &stored)
	assert.Equal(t, pool, stored)
}

func TestKeeper_GetPoolAmount_EmptyCoins(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	pool := k.GetPoolAmount(ctx)
	assert.Empty(t, pool)
}

func TestKeeper_GetPoolAmount_NonEmptyCoins(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	pool := sdk.NewCoins(sdk.NewInt64Coin("uatom", 100))
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.DepositsPoolStoreKey), cdc.MustMarshalBinaryBare(&pool))

	stored := k.GetPoolAmount(ctx)
	assert.Equal(t, pool, stored)
}
