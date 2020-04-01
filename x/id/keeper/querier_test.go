package keeper

import (
	"errors"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

// -----------------
// --- Identities
// -----------------

func Test_queryResolveIdentity_ExistingIdentity(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	var querier = NewQuerier(k)
	path := []string{types.QueryResolveDid, TestOwnerAddress.String()}
	actual, err := querier(ctx, path, request)
	require.Nil(t, err)

	expected, _ := codec.MarshalJSONIndent(cdc, ResolveIdentityResponse{
		Owner:       TestOwnerAddress,
		DidDocument: &TestDidDocument,
	})
	require.Equal(t, string(expected), string(actual))
}

func Test_queryResolveIdentity_nonExistentIdentity(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	var querier = NewQuerier(k)
	path := []string{types.QueryResolveDid, TestOwnerAddress.String()}
	actual, err := querier(ctx, path, request)
	require.Error(t, err)
	require.Nil(t, actual)
}

// -------------------
// --- Pairwise did
// -------------------

func Test_queryResolvePowerUpRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set(getDidPowerUpRequestStoreKey(TestDidPowerUpRequest.Proof), cdc.MustMarshalBinaryBare(&TestDidPowerUpRequest))

	var querier = NewQuerier(k)
	path := []string{types.QueryResolvePowerUpRequest, TestDidPowerUpRequest.Proof}
	actualBz, err := querier(ctx, path, request)

	var actual types.DidPowerUpRequest
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Nil(t, err)
	require.Equal(t, TestDidPowerUpRequest.Proof, actual.Proof)
	require.Equal(t, TestDidPowerUpRequest.Amount.String(), actual.Amount.String())
	require.Equal(t, TestDidPowerUpRequest.ID, actual.ID)
	require.Equal(t, TestDidPowerUpRequest.Status.Type, actual.Status.Type)
	require.Equal(t, TestDidPowerUpRequest.Status.Message, actual.Status.Message)
	require.Equal(t, TestDidPowerUpRequest.Claimant.String(), actual.Claimant.String())

}

func Test_queryResolvePowerUpRequest_NonExistingRequest(t *testing.T) {
	_, ctx, _, _, _, k := SetupTestInput()

	var querier = NewQuerier(k)
	path := []string{types.QueryResolvePowerUpRequest, ""}
	_, err := querier(ctx, path, request)

	require.Error(t, err)
	require.True(t, errors.Is(err, sdkErr.ErrUnknownRequest))
	require.Contains(t, err.Error(), "power-up request with id  not found")
}
