package keeper

import (
	"encoding/hex"
	"errors"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// ------------------
// --- Identities
// ------------------

func TestKeeper_CreateIdentity(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	err := k.SaveDidDocument(ctx, TestDidDocument)
	require.NoError(t, err)

	var stored types.DidDocument
	store := ctx.KVStore(k.storeKey)
	storedBz := store.Get(getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	require.Equal(t, TestDidDocument, stored)
}

func TestKeeper_EditIdentity(t *testing.T) {
	cdc, ctx, aK, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	updatedDocument := types.DidDocument{}
	err := copier.Copy(&updatedDocument, &TestDidDocument)
	require.NoError(t, err)

	account := aK.GetAccount(ctx, TestOwnerAddress)
	secp256k1key := account.GetPubKey().(secp256k1.PubKeySecp256k1)

	updatedDocument.PubKeys = types.PubKeys{
		types.PubKey{
			ID:         "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1",
			Type:       "Secp256k1VerificationKey2018",
			Controller: TestOwnerAddress,
			PublicKey:  hex.EncodeToString(secp256k1key[:]),
		},
	}

	err = k.SaveDidDocument(ctx, updatedDocument)
	require.NoError(t, err)

	var stored types.DidDocument
	storedBz := store.Get(getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	require.Equal(t, updatedDocument, stored)
	require.Len(t, stored.PubKeys, 1)
}

func TestKeeper_GetDidDocumentByOwner_ExistingDidDocument(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, err := k.GetDidDocumentByOwner(ctx, TestOwnerAddress)

	require.NoError(t, err)
	require.Equal(t, TestDidDocument, actual)
}

func TestKeeper_GetDidDocuments(t *testing.T) {
	cdc, ctx, aK, _, _, k := SetupTestInput()
	store := ctx.KVStore(k.storeKey)

	first := setupDidDocument(ctx, aK, "cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second := setupDidDocument(ctx, aK, "cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third := setupDidDocument(ctx, aK, "cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth := setupDidDocument(ctx, aK, "cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth := setupDidDocument(ctx, aK, "cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

	store.Set(getIdentityStoreKey(first.ID), cdc.MustMarshalBinaryBare(first))
	store.Set(getIdentityStoreKey(second.ID), cdc.MustMarshalBinaryBare(second))
	store.Set(getIdentityStoreKey(third.ID), cdc.MustMarshalBinaryBare(third))
	store.Set(getIdentityStoreKey(fourth.ID), cdc.MustMarshalBinaryBare(fourth))
	store.Set(getIdentityStoreKey(fifth.ID), cdc.MustMarshalBinaryBare(fifth))

	actual, err := k.GetDidDocuments(ctx)

	require.Nil(t, err)
	require.Equal(t, 5, len(actual))
	require.Contains(t, actual, first)
	require.Contains(t, actual, second)
	require.Contains(t, actual, third)
	require.Contains(t, actual, fourth)
	require.Contains(t, actual, fifth)
}

// ----------------------------
// --- Did power up requests
// ----------------------------

func TestKeeper_StorePowerUpRequest_NewRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	err := k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)
	require.Nil(t, err)

	var stored types.DidPowerUpRequest
	store := ctx.KVStore(k.storeKey)
	storedBz := store.Get(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	require.Equal(t, TestDidPowerUpRequest, stored)
}

func TestKeeper_StorePowerUpRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	err := k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "same proof")
}

func TestKeeper_GetPowerUpRequestByProof_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	_, err := k.GetPowerUpRequestByID(ctx, "")
	require.Error(t, err)
}

func TestKeeper_GetPowerUpRequestByProof_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	stored, err := k.GetPowerUpRequestByID(ctx, TestDidPowerUpRequest.Proof)
	require.NoError(t, err)
	require.Equal(t, TestDidPowerUpRequest, stored)
}

func TestKeeper_ChangePowerUpRequestStatus_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	status := types.RequestStatus{
		Type:    "status-type",
		Message: "status-messsge",
	}

	err := k.ChangePowerUpRequestStatus(ctx, "", status)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "proof")
}

func TestKeeper_ChangePowerUpRequestStatus_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	status := types.RequestStatus{
		Type:    "status-type",
		Message: "status-messsge",
	}

	err := k.ChangePowerUpRequestStatus(ctx, TestDidPowerUpRequest.Proof, status)
	require.Nil(t, err)

	expected := types.DidPowerUpRequest{
		Claimant:      TestDidPowerUpRequest.Claimant,
		Amount:        TestDidDepositRequest.Amount,
		Proof:         TestDidDepositRequest.Proof,
		EncryptionKey: TestDidDepositRequest.EncryptionKey,
		Status:        &status,
	}

	var stored types.DidPowerUpRequest
	storedBz := store.Get(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)
	require.Equal(t, expected, stored)
}

func TestKeeper_GetPowerUpRequests_EmptyList(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	didPowerUpRequests := k.GetPowerUpRequests(ctx)
	require.Empty(t, didPowerUpRequests)
}

func TestKeeper_GetPowerUpRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	didPowerUpRequests := k.GetPowerUpRequests(ctx)
	require.Equal(t, 1, len(didPowerUpRequests))
	require.Contains(t, didPowerUpRequests, TestDidPowerUpRequest)
}

func TestKeeper_GetPoolAmount_EmptyCoins(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	pool := k.GetPoolAmount(ctx)
	require.Empty(t, pool)
}

func TestKeeper_GetPoolAmount_NonEmptyCoins(t *testing.T) {
	_, ctx, _, bk, _, k := SetupTestInput()

	pool := sdk.NewCoins(sdk.NewInt64Coin("uatom", 100))
	k.supplyKeeper.SetSupply(ctx, supply.NewSupply(pool))
	_ = bk.SetCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName), pool)

	stored := k.GetPoolAmount(ctx)
	require.Equal(t, pool, stored)
}
