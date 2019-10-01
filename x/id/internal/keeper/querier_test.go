package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func Test_queryResolveIdentity_ExistingIdentity(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	var querier = NewQuerier(k)
	path := []string{"identities", TestOwnerAddress.String()}
	actual, err := querier(ctx, path, request)
	assert.Nil(t, err)

	expected, _ := codec.MarshalJSONIndent(cdc, ResolveIdentityResponse{
		Owner:       TestOwnerAddress,
		DidDocument: &TestDidDocument,
	})
	assert.Equal(t, string(expected), string(actual))
}

func Test_queryResolveIdentity_nonExistentIdentity(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	var querier = NewQuerier(k)
	path := []string{"identities", TestOwnerAddress.String()}
	actual, err := querier(ctx, path, request)
	assert.Nil(t, err)

	expected, _ := codec.MarshalJSONIndent(cdc, ResolveIdentityResponse{
		Owner:       TestOwnerAddress,
		DidDocument: nil,
	})
	assert.Equal(t, string(expected), string(actual))
}
