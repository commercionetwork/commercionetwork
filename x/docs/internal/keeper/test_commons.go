package keeper

import (
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an enviroment to test modules
func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// Store keys
	keyDocs := sdk.NewKVStoreKey("docs")
	keyGovernment := sdk.NewKVStoreKey("government")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keyDocs, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyGovernment, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	govk := government.NewKeeper(cdc, keyGovernment)
	dck := NewKeeper(keyDocs, govk, cdc)

	return cdc, ctx, dck
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}

// Testing variables

var TestingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestingSender2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var TestingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var TestingDocument = types.Document{
	Uuid:       "test-document-uuid",
	ContentUri: "https://example.com/document",
	Metadata: types.DocumentMetadata{
		ContentUri: "https://example.com/document/metadata",
		Schema: &types.DocumentMetadataSchema{
			Uri:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
	},
	Checksum: &types.DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
	Sender:     TestingSender,
	Recipients: ctypes.Addresses{TestingRecipient},
}

var TestingDocumentReceipt = types.DocumentReceipt{
	Uuid:         "testing-document-receipt-uuid",
	Sender:       TestingSender,
	Recipient:    TestingRecipient,
	TxHash:       "txHash",
	DocumentUuid: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}
