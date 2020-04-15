package keeper

import (
	"encoding/hex"
	"errors"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/id/types"
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
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	updatedDocument := types.DidDocument{}
	err := copier.Copy(&updatedDocument, &TestDidDocument)
	require.NoError(t, err)

	updatedDocument.PubKeys = append(updatedDocument.PubKeys, types.PubKey{
		ID:         updatedDocument.ID.String() + "#keys-3",
		Type:       "Secp256k1VerificationKey2018",
		Controller: updatedDocument.ID,
		PublicKey:  hex.EncodeToString([]byte("new key!")),
	})

	updatedDocument.Proof.SignatureValue = "K0scbXDo0D/vk1CoihZQ3gGn2ZjzQMi4zOZXVTgd7LQ5rZRBHAAYTNjJ/n0F24JiKa0b/bm52nKtzgy4DgP78w=="

	err = k.SaveDidDocument(ctx, updatedDocument)
	require.NoError(t, err)

	var stored types.DidDocument
	storedBz := store.Get(getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	require.Equal(t, updatedDocument, stored)
	require.Len(t, stored.PubKeys, 3)
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
	cdc, ctx, _, _, _, k := SetupTestInput()
	store := ctx.KVStore(k.storeKey)

	first := setupDidDocument()

	store.Set(getIdentityStoreKey(first.ID), cdc.MustMarshalBinaryBare(first))

	actual := k.GetDidDocuments(ctx)

	require.Equal(t, 1, len(actual))
	require.Contains(t, actual, first)
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
	storedBz := store.Get(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.ID))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	require.Equal(t, TestDidPowerUpRequest, stored)
}

func TestKeeper_StorePowerUpRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.ID), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	err := k.StorePowerUpRequest(ctx, TestDidPowerUpRequest)
	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "same id")
}

func TestKeeper_GetPowerUpRequestByProof_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	_, err := k.GetPowerUpRequestByID(ctx, "")
	require.Error(t, err)
}

func TestKeeper_GetPowerUpRequestByProof_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.ID), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	stored, err := k.GetPowerUpRequestByID(ctx, TestDidPowerUpRequest.ID)
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
	require.Contains(t, err.Error(), "id")
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
		Claimant: TestDidPowerUpRequest.Claimant,
		Amount:   TestDidPowerUpRequest.Amount,
		Proof:    TestDidPowerUpRequest.Proof,
		Status:   &status,
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

func TestKeeper_GetApprovedPowerUpRequests_EmptyList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	didPowerUpRequests := k.GetApprovedPowerUpRequests(ctx)

	require.Empty(t, didPowerUpRequests)
}

func TestKeeper_GetApprovedPowerUpRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)

	r := TestDidPowerUpRequest
	r.Status = &types.RequestStatus{
		Type:    types.StatusApproved,
		Message: "",
	}
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&r))

	didPowerUpRequests := k.GetApprovedPowerUpRequests(ctx)

	require.Equal(t, 1, len(didPowerUpRequests))
	require.Contains(t, didPowerUpRequests, r)
}

func TestKeeper_GetRejectedPowerUpRequests_EmptyList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	r := TestDidPowerUpRequest
	r.Status = &types.RequestStatus{
		Type: types.StatusRejected,
	}
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&r))

	didPowerUpRequests := k.GetRejectedPowerUpRequests(ctx)

	require.Contains(t, didPowerUpRequests, r)
}

func TestKeeper_GetRejectedPowerUpRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)

	r := TestDidPowerUpRequest
	r.Status = &types.RequestStatus{
		Type:    types.StatusRejected,
		Message: "",
	}
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&r))

	didPowerUpRequests := k.GetRejectedPowerUpRequests(ctx)

	require.Equal(t, 1, len(didPowerUpRequests))
	require.Contains(t, didPowerUpRequests, r)
}

func TestKeeper_GetPendingPowerUpRequests_EmptyList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)

	r := TestDidPowerUpRequest
	r.Status = nil
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&r))

	didPowerUpRequests := k.GetPendingPowerUpRequests(ctx)

	require.Contains(t, didPowerUpRequests, r)
}

func TestKeeper_GetPendingPowerUpRequests_ExistingList(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)

	r := TestDidPowerUpRequest
	r.Status = &types.RequestStatus{
		Type:    types.StatusApproved,
		Message: "",
	}
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&r))

	didPowerUpRequests := k.GetPendingPowerUpRequests(ctx)

	require.Equal(t, 0, len(didPowerUpRequests))
	require.NotContains(t, didPowerUpRequests, r)
}
