package commerciodocs

import (
	"commercio-network/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//testing add readers
func TestKeeper_addReaderForDocument(t *testing.T) {

	var readers = []types.Did{"reader", "reader2"}

	store := input.ctx.KVStore(input.docsKeeper.readersStoreKey)
	store.Set([]byte(reference), input.cdc.MustMarshalBinaryBare(&readers))

	currentLength := len(store.Get([]byte(reference)))

	input.docsKeeper.addReaderForDocument(input.ctx, ownerIdentity, reference)

	afterOpLength := len(store.Get([]byte(reference)))

	if afterOpLength < currentLength {
		t.Errorf("afterOpLength should be greater than currentLength")
	}
}

func TestKeeper_StoreDocument(t *testing.T) {

	ownerStore := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	ownerStore.Set([]byte(reference), []byte(owner))

	metadataStore := input.ctx.KVStore(input.docsKeeper.metadataStoreKey)

	currentLength := len(metadataStore.Get([]byte(reference)))

	input.docsKeeper.StoreDocument(input.ctx, owner, ownerIdentity, reference, metadata)

	afterOpLength := len(metadataStore.Get([]byte(reference)))

	if afterOpLength < currentLength {
		t.Errorf("after operation length should be greater than current length")
	}
}

//Given reference has an owner
func TestKeeper_HasOwner_True(t *testing.T) {

	ass := assert.New(t)

	store := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	store.Set([]byte(reference), []byte(owner))

	result := input.docsKeeper.HasOwner(input.ctx, reference)

	ass.True(result)
}

//Given reference hasn't got an owner
func TestKeeper_HasOwner_False(t *testing.T) {
	ass := assert.New(t)

	result := input.docsKeeper.HasOwner(input.ctx, reference)

	ass.False(result)
}

//Given owner is the owner of doc reference
func TestKeeper_IsOwner_True(t *testing.T) {

	ass := assert.New(t)

	store := input.ctx.KVStore(input.docsKeeper.ownersStoreKey)
	store.Set([]byte(reference), []byte(owner))

	res := input.docsKeeper.IsOwner(input.ctx, owner, reference)

	ass.True(res)
}

//Given owner isnt the owner of the doc reference
func TestKeeper_IsOwner_False(t *testing.T) {

	ass := assert.New(t)

	res := input.docsKeeper.IsOwner(input.ctx, owner, reference)

	ass.False(res)
}

func TestKeeper_GetMetadata_OfExistentDocument(t *testing.T) {

	ass := assert.New(t)

	metadataStore := input.ctx.KVStore(input.docsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(reference), []byte(metadata))

	result := input.docsKeeper.GetMetadata(input.ctx, reference)

	ass.Equal(metadata, result)

}

func TestKeeper_GetMetadata_OfNonExistentDocument(t *testing.T) {

	ass := assert.New(t)

	result := input.docsKeeper.GetMetadata(input.ctx, reference)

	ass.Equal("", result)
}

func TestKeeper_CanReadDocument_True(t *testing.T) {

	ass := assert.New(t)

	readers := []types.Did{ownerIdentity}

	readerStore := input.ctx.KVStore(input.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(reference), input.cdc.MustMarshalBinaryBare(&readers))

	result := input.docsKeeper.CanReadDocument(input.ctx, ownerIdentity, reference)

	ass.True(result)
}

func TestKeeper_CanReadDocument_False(t *testing.T) {

	ass := assert.New(t)

	result := input.docsKeeper.CanReadDocument(input.ctx, ownerIdentity, reference)

	ass.False(result)
}

func TestKeeper_GetAuthorizedReaders(t *testing.T) {
	var readers = []types.Did{"reader", "reader2"}

	store := input.ctx.KVStore(input.docsKeeper.readersStoreKey)
	store.Set([]byte(reference), input.cdc.MustMarshalBinaryBare(&readers))

	res := input.docsKeeper.GetAuthorizedReaders(input.ctx, reference)

	assert.Equal(t, readers, res)
}

func TestKeeper_ShareDocument_SenderAuthorizedToShare(t *testing.T) {

	var readers = []types.Did{ownerIdentity}

	readerStore := input.ctx.KVStore(input.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(reference), input.cdc.MustMarshalBinaryBare(&readers))

	result := input.docsKeeper.ShareDocument(input.ctx, reference, ownerIdentity, recipient)

	assert.Nil(t, result)
}

func TestKeeper_ShareDocument_SenderUnauthorizedToShare(t *testing.T) {

	error := sdk.ErrUnauthorized(fmt.Sprintf("The sender with address %s doesnt have the rights on this document", ownerIdentity))

	result := input.docsKeeper.ShareDocument(input.ctx, reference, ownerIdentity, recipient)

	assert.NotNil(t, result)

	assert.Equal(t, error, result)
}
