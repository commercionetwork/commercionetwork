package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

// Testing variables
var testingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var anotherTestingSender, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var testingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var anotherTestingRecipient, _ = sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")

const anotherValidDocumentUUID = "49c981c2-a09e-47d2-8814-9373ff64abae"
const documentReceiptUUID = "32c82ee4-c71d-4890-9680-4db7a3dbed41"
const anotherDocumentReceiptUUID = "4c24eda0-6c06-476b-99ab-a05ea6f3d14f"

var testingDocument = types.Document{
	UUID:       "d83422c6-6e79-4a99-9767-fcae46dfa371",
	ContentURI: "https://example.com/document",
	Metadata: &types.DocumentMetadata{
		ContentURI: "https://example.com/document/metadata",
		Schema: &types.DocumentMetadataSchema{
			URI:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
	},
	Checksum: &types.DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
	Sender:     testingSender.String(),
	Recipients: append([]string{}, testingRecipient.String(), anotherTestingSender.String()),
}

var testingDocumentReceipt = types.DocumentReceipt{
	UUID:         documentReceiptUUID,
	Sender:       testingSender.String(),
	Recipient:    testingRecipient.String(),
	TxHash:       "txHash",
	DocumentUUID: testingDocument.UUID,
	Proof:        "proof",
}

var testingDocumentReceiptNoDoc = types.DocumentReceipt{
	UUID:         anotherDocumentReceiptUUID,
	Sender:       testingSender.String(),
	Recipient:    testingRecipient.String(),
	TxHash:       "txHash",
	DocumentUUID: anotherValidDocumentUUID,
	Proof:        "proof",
}

func createNDocument(keeper *Keeper, ctx sdk.Context, n int) []*types.Document {
	items := []*types.Document{}
	for i := 0; i < n; i++ {
		item := &types.Document{
			Sender:     testingSender.String(),
			Recipients: []string{testingRecipient.String()},
			UUID:       uuid.NewV4().String(),
		}
		items = append(items, item)

		_ = keeper.SaveDocument(ctx, *items[i])
	}
	return items
}

func createNDocumentReceipt(keeper *Keeper, ctx sdk.Context, n int) []*types.DocumentReceipt {
	docs := createNDocument(keeper, ctx, n)

	items := []*types.DocumentReceipt{}
	for i := range docs {
		item := &types.DocumentReceipt{
			Sender:       docs[i].Recipients[0],
			DocumentUUID: docs[i].UUID,
			Recipient:    docs[i].Sender,
			UUID:         uuid.NewV4().String(),
		}
		items = append(items, item)

		_ = keeper.SaveReceipt(ctx, *items[i])
	}
	return items
}

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(codec.NewProtoCodec(registry), storeKey, memStoreKey)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return keeper, ctx
}
