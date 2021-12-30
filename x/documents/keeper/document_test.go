package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func createNDocument(keeper *Keeper, ctx sdk.Context, n int) []types.Document {
	items := make([]types.Document, n)
	for i := range items {
		items[i].Sender = testingSender.String()
		items[i].UUID = uuid.NewV4().String()
		//_ = keeper.AppendDocument(ctx, items[i])
		_ = keeper.SaveDocument(ctx, items[i])
	}
	return items
}

func TestDocumentGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	for _, item := range items {
		actual, _ := keeper.GetDocumentByID(ctx, item.UUID)
		assert.Equal(t, item, actual)
	}
}

func TestDocumentExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasDocument(ctx, item.UUID))
	}
}

func TestDocumentGetAll(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	actual := keeper.GetAllDocument(ctx)
	//assert.Equal(t, items, keeper.GetAllDocument(ctx))
	assert.Len(t, actual, len(items))

	for _, item := range items {
		assert.Contains(t, actual, item)
	}
}
