package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetGovernmentAddress_NonExisting(t *testing.T) {
	input := setupTestInput()

	err := input.Keeper.SetGovernmentAddress(input.Ctx, TestAddress)
	assert.Nil(t, err)

	store := input.Ctx.KVStore(input.Keeper.StoreKey)
	stored := sdk.AccAddress(store.Get([]byte(types.GovernmentStoreKey)))
	assert.Equal(t, TestAddress, stored)
}

func TestKeeper_GetGovernmentAddress(t *testing.T) {
	input := setupTestInput()
	store := input.Ctx.KVStore(input.Keeper.StoreKey)
	store.Set([]byte(types.GovernmentStoreKey), TestAddress)

	actual := input.Keeper.GetGovernmentAddress(input.Ctx)

	assert.Equal(t, TestAddress, actual)
}
