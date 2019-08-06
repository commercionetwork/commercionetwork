package keeper

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//testing add readers
func TestKeeper_addReaderForDocument(t *testing.T) {

	var readers = []types.Did{"reader", "reader2"}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.readersStoreKey)
	store.Set([]byte(TestReference), TestUtils.Cdc.MustMarshalBinaryBare(&readers))

	currentLength := len(store.Get([]byte(TestReference)))

	TestUtils.DocsKeeper.addReaderForDocument(TestUtils.Ctx, TestOwnerIdentity, TestReference)

	afterOpLength := len(store.Get([]byte(TestReference)))

	if afterOpLength < currentLength {
		t.Errorf("afterOpLength should be greater than currentLength")
	}
}

func TestKeeper_StoreDocument(t *testing.T) {

	ownerStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.ownersStoreKey)
	ownerStore.Set([]byte(TestReference), []byte(TestOwner))

	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.metadataStoreKey)

	currentLength := len(metadataStore.Get([]byte(TestReference)))

	TestUtils.DocsKeeper.StoreDocument(TestUtils.Ctx, TestOwner, TestOwnerIdentity, TestReference, TestMetadata)

	afterOpLength := len(metadataStore.Get([]byte(TestReference)))

	if afterOpLength < currentLength {
		t.Errorf("after operation length should be greater than current length")
	}
}

//Given TestReference has an TestOwner
func TestKeeper_HasOwner_True(t *testing.T) {

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.ownersStoreKey)
	store.Set([]byte(TestReference), []byte(TestOwner))

	result := TestUtils.DocsKeeper.HasOwner(TestUtils.Ctx, TestReference)

	assert.True(t, result)
}

//Given TestReference hasn't got an TestOwner
func TestKeeper_HasOwner_False(t *testing.T) {

	reference := "reff"

	result := TestUtils.DocsKeeper.HasOwner(TestUtils.Ctx, reference)

	assert.False(t, result)
}

//Given TestOwner is the TestOwner of doc TestReference
func TestKeeper_IsOwner_True(t *testing.T) {

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.ownersStoreKey)
	store.Set([]byte(TestReference), []byte(TestOwner))

	res := TestUtils.DocsKeeper.IsOwner(TestUtils.Ctx, TestOwner, TestReference)

	assert.True(t, res)
}

//Given TestOwner isnt the TestOwner of the doc TestReference
func TestKeeper_IsOwner_False(t *testing.T) {

	reference := "reff"
	res := TestUtils.DocsKeeper.IsOwner(TestUtils.Ctx, TestOwner, reference)

	assert.False(t, res)
}

func TestKeeper_GetMetadata_OfExistentDocument(t *testing.T) {

	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(TestReference), []byte(TestMetadata))

	result := TestUtils.DocsKeeper.GetMetadata(TestUtils.Ctx, TestReference)

	assert.Equal(t, TestMetadata, result)
}

func TestKeeper_GetMetadata_OfNonExistentDocument(t *testing.T) {

	reference := "reff"

	result := TestUtils.DocsKeeper.GetMetadata(TestUtils.Ctx, reference)

	assert.Equal(t, "", result)
}

func TestKeeper_CanReadDocument_True(t *testing.T) {

	readers := []types.Did{TestOwnerIdentity}

	readerStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.readersStoreKey)
	readerStore.Set([]byte(TestReference), TestUtils.Cdc.MustMarshalBinaryBare(&readers))

	result := TestUtils.DocsKeeper.CanReadDocument(TestUtils.Ctx, TestOwnerIdentity, TestReference)

	assert.True(t, result)
}

func TestKeeper_CanReadDocument_False(t *testing.T) {

	reference := "reff"

	result := TestUtils.DocsKeeper.CanReadDocument(TestUtils.Ctx, TestOwnerIdentity, reference)

	assert.False(t, result)
}

func TestKeeper_GetAuthorizedReaders(t *testing.T) {
	var readers = []types.Did{"reader", "reader2"}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.readersStoreKey)
	store.Set([]byte(TestReference), TestUtils.Cdc.MustMarshalBinaryBare(&readers))

	res := TestUtils.DocsKeeper.GetAuthorizedReaders(TestUtils.Ctx, TestReference)

	assert.Equal(t, readers, res)
}

func TestKeeper_ShareDocument_SenderAuthorizedToShare(t *testing.T) {

	var readers = []types.Did{TestOwnerIdentity}

	readerStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.readersStoreKey)
	readerStore.Set([]byte(TestReference), TestUtils.Cdc.MustMarshalBinaryBare(&readers))

	result := TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestReference, TestOwnerIdentity, TestRecipient)

	assert.Nil(t, result)
}

func TestKeeper_ShareDocument_SenderUnauthorizedToShare(t *testing.T) {

	ownerIdentity := types.Did("notOwner")
	error := sdk.ErrUnauthorized(fmt.Sprintf("The sender with TestAddress %s doesnt have the rights on this document", ownerIdentity))

	result := TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestReference, ownerIdentity, TestRecipient)

	assert.NotNil(t, result)

	assert.Equal(t, error, result)
}
