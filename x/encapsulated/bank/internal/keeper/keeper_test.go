package keeper

import (
	"testing"

	cmtypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestAddress, _ = sdk.AccAddressFromBech32("cosmos1mlrqrdrxs50z972h32x9w8x3lta7hkms0hxraq")
var TestAddress2, _ = sdk.AccAddressFromBech32("cosmos1dy7v2cfgunggrul5sqyj9sunxpyyhtcx94hlnn")

func TestKeeper_AddBlockedAddresses_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	k.AddBlockedAddresses(ctx, TestAddress)

	var addresses []sdk.AccAddress
	store := ctx.KVStore(k.storeKey)
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.BlockedAddressesStoreKey)), &addresses)

	assert.Len(t, addresses, 1)
	assert.Contains(t, addresses, TestAddress)
}

func TestKeeper_AddBlockedAddresses_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	existing := []sdk.AccAddress{TestAddress}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.BlockedAddressesStoreKey), cdc.MustMarshalBinaryBare(&existing))

	k.AddBlockedAddresses(ctx, TestAddress2)

	var addresses []sdk.AccAddress
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.BlockedAddressesStoreKey)), &addresses)

	assert.Len(t, addresses, 2)
	assert.Contains(t, addresses, TestAddress)
	assert.Contains(t, addresses, TestAddress2)
}

func TestKeeper_GetBlockedAddresses_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()
	stored := k.GetBlockedAddresses(ctx)
	assert.Len(t, stored, 0)
}

func TestKeeper_GetBlockedAddresses_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	existing := cmtypes.Addresses{TestAddress, TestAddress2}
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.BlockedAddressesStoreKey), cdc.MustMarshalBinaryBare(&existing))

	stored := k.GetBlockedAddresses(ctx)
	assert.Len(t, stored, 2)
	assert.Equal(t, existing, stored)
}
