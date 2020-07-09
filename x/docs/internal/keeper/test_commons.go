package keeper

import (
	"fmt"
	"sync"
	"time"

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
func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, keeper Keeper, idKeeper idkeeper.Keeper, membershipsKeeper memberships.Keeper, bankKeeper bank.Keeper, idStoreKey *sdk.KVStoreKey) {

	memDB := db.NewMemDB()
	cdc = testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	supKey := sdk.NewKVStoreKey(supply.StoreKey)
	idsk := sdk.NewKVStoreKey(idtypes.StoreKey)
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	memk := sdk.NewKVStoreKey(memberships.StoreKey)
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
	ms.MountStoreWithDB(idsk, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyGovernment, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(memk, sdk.StoreTypeIAVL, memDB)
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

	memkeeper := memberships.NewKeeper(cdc, memk, sk, govk, ak)

	dck := NewKeeper(keyDocs, govk, bk, idk, memkeeper, cdc)

	return cdc, ctx, dck, idk, memkeeper, bk, idsk
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

func setupDidDocument() idtypes.DidDocument {
	var testZone, _ = time.LoadLocation("UTC")
	var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)
	var testOwnerAddress, _ = sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")

	return idtypes.DidDocument{
		Context: idtypes.ContextDidV1,
		ID:      testOwnerAddress,
		Proof: &idtypes.Proof{
			Type:               idtypes.KeyTypeSecp256k12019,
			Created:            testTime,
			ProofPurpose:       idtypes.ProofPurposeAuthentication,
			Controller:         testOwnerAddress.String(),
			SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
			VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
		},
		PubKeys: idtypes.PubKeys{
			idtypes.PubKey{
				ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
				Type:       "RsaVerificationKey2018",
				Controller: testOwnerAddress,
				PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
			},
			idtypes.PubKey{
				ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
				Type:       "RsaSignatureKey2018",
				Controller: testOwnerAddress,
				PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
			},
		},
	}
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
