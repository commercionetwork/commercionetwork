package documents

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/commercionetwork/commercionetwork/x/documents/keeper"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	storemetrics "cosmossdk.io/store/metrics"
)

func setupKeeper(t testing.TB) (*keeper.Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := keeper.NewKeeper(codec.NewProtoCodec(registry), storeKey, memStoreKey)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return keeper, ctx
}

func TestInitGenesis(t *testing.T) {
	tests := []struct {
		name      string
		docs      []*types.Document
		receipts  []*types.DocumentReceipt
		wantPanic bool
	}{
		{
			name:      "invalid document",
			docs:      []*types.Document{&types.InvalidDocument},
			receipts:  []*types.DocumentReceipt{},
			wantPanic: true,
		},
		{
			name:      "duplicate document",
			docs:      []*types.Document{&types.ValidDocument, &types.ValidDocument},
			receipts:  []*types.DocumentReceipt{},
			wantPanic: true,
		},
		{
			name:      "invalid receipt",
			docs:      []*types.Document{&types.ValidDocument},
			receipts:  []*types.DocumentReceipt{&types.InvalidDocumentReceipt},
			wantPanic: true,
		},
		{
			name:      "duplicate receipt",
			docs:      []*types.Document{&types.ValidDocument},
			receipts:  []*types.DocumentReceipt{&types.ValidDocumentReceiptRecipient1, &types.ValidDocumentReceiptRecipient1},
			wantPanic: true,
		},
		{
			name:     "empty",
			docs:     []*types.Document{},
			receipts: []*types.DocumentReceipt{},
		},
		{
			name:     "empty receipts",
			docs:     []*types.Document{&types.ValidDocument},
			receipts: []*types.DocumentReceipt{},
		},
		{
			name: "one receipt",
			docs: []*types.Document{&types.ValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
			},
		},
		{
			name: "two receipts",
			docs: []*types.Document{&types.ValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
				&types.ValidDocumentReceiptRecipient2,
			},
		},
		{
			name: "two documents, two receipts",
			docs: []*types.Document{&types.ValidDocument, &types.AnotherValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
				&types.ValidDocumentReceiptRecipient2,
			},
		},
		{
			name: "two documents, three receipts",
			docs: []*types.Document{&types.ValidDocument, &types.AnotherValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
				&types.ValidDocumentReceiptRecipient2,
				&types.AnotherValidDocumentReceipt,
			},
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			defer func() { recover() }()

			InitGenesis(ctx, *keeper, types.GenesisState{
				Documents: tt.docs,
				Receipts:  tt.receipts,
			})

			defer func() {
				if tt.wantPanic {
					t.Error("should have panicked")
				}
			}()

			got := ExportGenesis(ctx, *keeper)
			require.ElementsMatch(t, tt.docs, got.Documents)
			require.ElementsMatch(t, tt.receipts, got.Receipts)
		})
	}
}

func TestExportGenesis(t *testing.T) {
	tests := []struct {
		name      string
		recipient string
		docs      []*types.Document
		receipts  []*types.DocumentReceipt
	}{
		{
			name:      "empty",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []*types.Document{},
			receipts:  []*types.DocumentReceipt{},
		},
		{
			name:      "empty receipts",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []*types.Document{&types.ValidDocument},
			receipts:  []*types.DocumentReceipt{},
		},
		{
			name:      "one receipt",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []*types.Document{&types.ValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
			},
		},
		{
			name:      "two receipts",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []*types.Document{&types.ValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
				&types.ValidDocumentReceiptRecipient2,
			},
		},
		{
			name:      "two documents, two receipts",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []*types.Document{&types.ValidDocument, &types.AnotherValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
				&types.ValidDocumentReceiptRecipient2,
			},
		},
		{
			name:      "two documents, three receipts",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []*types.Document{&types.ValidDocument, &types.AnotherValidDocument},
			receipts: []*types.DocumentReceipt{
				&types.ValidDocumentReceiptRecipient1,
				&types.ValidDocumentReceiptRecipient2,
				&types.AnotherValidDocumentReceipt,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.docs {
				keeper.SaveDocument(ctx, *document)
			}

			for _, receipt := range tt.receipts {
				keeper.SaveReceipt(ctx, *receipt)
			}

			got := ExportGenesis(ctx, *keeper)
			require.ElementsMatch(t, tt.docs, got.Documents)
			require.ElementsMatch(t, tt.receipts, got.Receipts)
		})
	}
}
