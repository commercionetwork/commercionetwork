package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

func TestKeeper_SetGovernmentAddress_NonExisting(t *testing.T) {
	_, ctx, k := SetupTestInput(false)

	err := k.SetGovernmentAddress(ctx, governmentTestAddress)
	require.Nil(t, err)

	store := ctx.KVStore(k.StoreKey)
	stored := sdk.AccAddress(store.Get([]byte(types.GovernmentStoreKey)))
	require.Equal(t, governmentTestAddress, stored)
}

func TestKeeper_GetGovernmentAddress(t *testing.T) {
	_, ctx, k := SetupTestInput(true)
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.GovernmentStoreKey), governmentTestAddress)

	actual := k.GetGovernmentAddress(ctx)

	require.Equal(t, governmentTestAddress, actual)
}

func TestKeeper_SetTumblerAddress_NonExisting(t *testing.T) {
	_, ctx, k := SetupTestInput(false)

	err := k.SetTumblerAddress(ctx, governmentTestAddress)
	require.Nil(t, err)

	store := ctx.KVStore(k.StoreKey)
	stored := sdk.AccAddress(store.Get([]byte(types.TumblerStoreKey)))
	require.Equal(t, governmentTestAddress, stored)
}

func TestKeeper_GetTumblerAddress(t *testing.T) {
	_, ctx, k := SetupTestInput(true)
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.TumblerStoreKey), governmentTestAddress)

	actual := k.GetTumblerAddress(ctx)

	require.Equal(t, governmentTestAddress, actual)
}
