package keeper

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func SetupTestInput() (*codec.Codec, sdk.Context, auth.AccountKeeper, Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// CommercioID
	storeKey := sdk.NewKVStoreKey("commercioid")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	idk := NewKeeper(cdc, storeKey, ak)

	// Setup the Did Document
	TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	TestDidDocument = setupDidDocument(ctx, ak, "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

	return cdc, ctx, ak, idk

}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.PubKeyAminoName, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.PrivKeyAminoName, nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{}, ed25519.PubKeyAminoName, nil)
	cdc.RegisterConcrete(ed25519.PrivKeyEd25519{}, ed25519.PrivKeyAminoName, nil)

	auth.RegisterCodec(cdc)
	types.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

func setupDidDocument(ctx sdk.Context, ak auth.AccountKeeper, bech32Address string) types.DidDocument {
	// Create a public key
	var secp256k1Key secp256k1.PubKeySecp256k1
	bz, _ := hex.DecodeString("02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71")
	copy(secp256k1Key[:], bz)

	// Create the owner account
	address, _ := sdk.AccAddressFromBech32(bech32Address)
	account := ak.NewAccountWithAddress(ctx, address)
	_ = account.SetPubKey(secp256k1Key)
	ak.SetAccount(ctx, account)

	testZone, _ := time.LoadLocation("UTC")
	testTime := time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)

	return types.DidDocument{
		Context: "https://www.w3.org/2019/did/v1",
		Id:      address,
		Authentication: []string{
			fmt.Sprintf("%s#keys-1", address),
		},
		Proof: types.Proof{
			Type:           "LinkedDataSignature2015",
			Created:        testTime,
			Creator:        fmt.Sprintf("%s#keys-1", address),
			SignatureValue: "QNB13Y7Q9...1tzjn4w==",
		},
		PubKeys: types.PubKeys{
			types.PubKey{
				Id:           fmt.Sprintf("%s#keys-1", address),
				Type:         "Secp256k1VerificationKey2018",
				Controller:   address,
				PublicKeyHex: hex.EncodeToString(account.GetPubKey().Bytes()),
			},
			types.PubKey{
				Id:           fmt.Sprintf("%s#keys-2", address),
				Type:         "RsaVerificationKey2018",
				Controller:   address,
				PublicKeyHex: "04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a",
			},
		},
	}
}

// Test variables
var TestOwnerAddress sdk.AccAddress
var TestDidDocument types.DidDocument
