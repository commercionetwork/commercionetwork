package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

// -----------------
// --- Identities
// -----------------

func Test_queryResolveIdentity_ExistingIdentity(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getIdentityStoreKey(TestOwnerAddress), cdc.MustMarshalBinaryBare(TestDidDocument))

	var querier = NewQuerier(k)
	path := []string{types.QueryResolveDid, TestOwnerAddress.String()}
	actual, err := querier(ctx, path, request)
	assert.Nil(t, err)

	expected, _ := codec.MarshalJSONIndent(cdc, ResolveIdentityResponse{
		Owner:       TestOwnerAddress,
		DidDocument: &TestDidDocument,
	})
	assert.Equal(t, string(expected), string(actual))
}

func Test_queryResolveIdentity_nonExistentIdentity(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	var querier = NewQuerier(k)
	path := []string{types.QueryResolveDid, TestOwnerAddress.String()}
	actual, err := querier(ctx, path, request)
	assert.Nil(t, err)

	expected, _ := codec.MarshalJSONIndent(cdc, ResolveIdentityResponse{
		Owner:       TestOwnerAddress,
		DidDocument: nil,
	})
	assert.Equal(t, string(expected), string(actual))
}

// -------------------
// --- Pairwise did
// -------------------

func Test_queryResolveDepositRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)

	depositRequest := types.DidDepositRequest{
		Recipient:     requestRecipient,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
		FromAddress:   requestSender,
	}
	store.Set(k.getDepositRequestStoreKey(depositRequest.Proof), cdc.MustMarshalBinaryBare(&depositRequest))

	var querier = NewQuerier(k)
	path := []string{types.QueryResolveDepositRequest, depositRequest.Proof}
	actualBz, err := querier(ctx, path, request)

	var actual types.DidDepositRequest
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Nil(t, err)
	assert.Equal(t, depositRequest, actual)
}

func Test_queryResolveDepositRequest_NonExistingRequest(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	var querier = NewQuerier(k)
	path := []string{types.QueryResolveDepositRequest, ""}
	_, err := querier(ctx, path, request)

	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "proof")
}

func Test_queryResolvePowerUpRequest_ExistingRequest(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)

	powerUpRequest := types.DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "proof",
		EncryptionKey: "encryption_key",
	}
	store.Set(k.getDidPowerUpRequestStoreKey(powerUpRequest.Proof), cdc.MustMarshalBinaryBare(&powerUpRequest))

	var querier = NewQuerier(k)
	path := []string{types.QueryResolvePowerUpRequest, powerUpRequest.Proof}
	actualBz, err := querier(ctx, path, request)

	var actual types.DidPowerUpRequest
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Nil(t, err)
	assert.Equal(t, powerUpRequest, actual)
}

func Test_queryResolvePowerUpRequest_NonExistingRequest(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	var querier = NewQuerier(k)
	path := []string{types.QueryResolvePowerUpRequest, ""}
	_, err := querier(ctx, path, request)

	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "proof")
}
