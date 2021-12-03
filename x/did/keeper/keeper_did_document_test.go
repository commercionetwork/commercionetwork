package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIdentityGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNIdentityNew(keeper, ctx, 10)
	for _, item := range items {
		a, err := keeper.GetDdoByOwner(ctx, sdk.AccAddress(item.ID))
		require.NoError(t, err)
		assert.Equal(t, item, a)
	}
}

func TestNewDocumentExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNIdentityNew(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasIdentity(ctx, item.ID))
	}
}

func createNIdentityNew(keeper *Keeper, ctx sdk.Context, n int) []types.DidDocument {
	items := make([]types.DidDocument, n)
	for i := range items {
		_, _, addr := testdata.KeyTestPubAddr()
		items[i].ID = string(addr)
		_ = keeper.AppendDid(ctx, items[i])
	}
	return items
}
