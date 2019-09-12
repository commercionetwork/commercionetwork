package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetGovernmentAddress_NonExisting(t *testing.T) {
	_, ctx, k := SetupTestInput()

	err := k.SetGovernmentAddress(ctx, TestAddress)
	assert.Nil(t, err)

	store := ctx.KVStore(k.StoreKey)
	stored := sdk.AccAddress(store.Get([]byte(types.GovernmentStoreKey)))
	assert.Equal(t, TestAddress, stored)
}

func TestKeeper_GetGovernmentAddress(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.GovernmentStoreKey), TestAddress)

	actual := k.GetGovernmentAddress(ctx)

	assert.Equal(t, TestAddress, actual)
}
