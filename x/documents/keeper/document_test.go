package keeper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	for _, item := range items {
		actual, _ := keeper.GetDocumentByID(ctx, item.UUID)
		assert.Equal(t, *item, actual)
	}
}

func TestDocumentExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasDocument(ctx, item.UUID))
	}
}
