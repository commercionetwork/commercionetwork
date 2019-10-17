package keeper

import (
	"encoding/hex"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_CreateIdentity(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	err := k.SaveDidDocument(ctx, TestDidDocument)
	assert.NoError(t, err)

	var stored types.DidDocument
	store := ctx.KVStore(k.storeKey)
	storedBz := store.Get(k.getIdentityStoreKey(TestOwnerAddress))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, TestDidDocument, stored)
}

func TestKeeper_EditIdentity(t *testing.T) {
	cdc, ctx, ak, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	updatedDocument := types.DidDocument{}
	err := copier.Copy(&updatedDocument, &TestDidDocument)
	assert.NoError(t, err)

	account := ak.GetAccount(ctx, TestOwnerAddress)
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
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	actual, found := k.GetDidDocumentByOwner(ctx, TestOwnerAddress)

	assert.True(t, found)
	assert.Equal(t, TestDidDocument, actual)
}

func TestKeeper_GetDidDocuments(t *testing.T) {
	cdc, ctx, ak, k := SetupTestInput()
	store := ctx.KVStore(k.storeKey)

	first := setupDidDocument(ctx, ak, "cosmos18xffcd029jn3thr0wwxah6gjdldr3kchvydkuj")
	second := setupDidDocument(ctx, ak, "cosmos18t0e6fevehhjv682gkxpchvmnl7z7ue4t4w0nd")
	third := setupDidDocument(ctx, ak, "cosmos1zt9etyl07asvf32g0d7ddjanres2qt9cr0fek6")
	fourth := setupDidDocument(ctx, ak, "cosmos177ap6yqt87znxmep5l7vdaac59uxyn582kv0gl")
	fifth := setupDidDocument(ctx, ak, "cosmos1ajv8j3e0ud2uduzdqmxfcvwm3nwdgr447yvu5m")

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
