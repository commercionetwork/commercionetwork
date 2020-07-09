package keeper

import (
	"sync"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	idkeeper "github.com/commercionetwork/commercionetwork/x/id/keeper"
	idtypes "github.com/commercionetwork/commercionetwork/x/id/types"
	"github.com/commercionetwork/commercionetwork/x/memberships"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

var (
	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32MainPrefix = "did:com:"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic

	configSealOnce sync.Once
)

//This function create an environment to test modules
func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, keeper Keeper, idKeeper idkeeper.Keeper, membershipsKeeper memberships.Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	supKey := sdk.NewKVStoreKey(supply.StoreKey)
	idsk := sdk.NewKVStoreKey(idtypes.StoreKey)
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	configSealOnce.Do(func() {
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
		config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
		config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
		config.Seal()
	})

	// Store keys
	keyDocs := sdk.NewKVStoreKey("docs")
	keyGovernment := sdk.NewKVStoreKey("government")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(supKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keyDocs, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyGovernment, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), map[string]bool{
		types.ModuleName: false,
	})
	maccPerms := map[string][]string{
		types.ModuleName:       nil,
		memberships.ModuleName: {supply.Burner},
		idtypes.ModuleName:     nil,
	}
	sk := supply.NewKeeper(cdc, supKey, ak, bk, maccPerms)
	govk := government.NewKeeper(cdc, keyGovernment)

	testGovernment, err := sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	if err != nil {
		panic(err)
	}

	err = govk.SetTumblerAddress(ctx, testGovernment)
	if err != nil {
		panic(err)
	}
	idk := idkeeper.NewKeeper(cdc, idsk, ak, sk)

	// Set initial supply
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 1))))

	memk := sdk.NewKVStoreKey(memberships.StoreKey)
	memkeeper := memberships.NewKeeper(cdc, memk, sk, govk, ak)

	dck := NewKeeper(keyDocs, govk, bk, idk, memkeeper, cdc)

	return cdc, ctx, dck, idk, memkeeper
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Testing variables

var TestingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestingSender2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var TestingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var TestingDocument = types.Document{
	UUID:       "test-document-uuid",
	ContentURI: "https://example.com/document",
	Metadata: types.DocumentMetadata{
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
	Sender:     TestingSender,
	Recipients: ctypes.Addresses{TestingRecipient},
}

var TestingDocumentReceipt = types.DocumentReceipt{
	UUID:         "testing-document-receipt-uuid",
	Sender:       TestingSender,
	Recipient:    TestingRecipient,
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}
