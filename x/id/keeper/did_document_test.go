package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func createNIdentity(keeper *Keeper, ctx sdk.Context, n int) []types.DidDocument {
	items := make([]types.DidDocument, n)
	for i := range items {
		_, _, addr := testdata.KeyTestPubAddr()
		items[i].ID = string(addr)
		_ = keeper.AppendId(ctx, items[i])
	}
	return items
}

/*
func TestIdentityGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNIdentity(keeper, ctx, 10)
	for _, item := range items {
		assert.Equal(t, item, keeper.GetIdentity(ctx, item.ID))
	}
}*/

func TestDocumentExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNIdentity(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasIdentity(ctx, item.ID))
	}
}
