package keeper

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/bech32"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func SetupTestInput() (*codec.Codec, sdk.Context, auth.AccountKeeper, bank.Keeper, government.Keeper, Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey,
		params.StoreKey,
		supply.StoreKey,
		government.StoreKey,
		types.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tkey := range tkeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, memDB)
	}
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	ak := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), map[string]bool{
		types.ModuleName: false,
	})
	maccPerms := map[string][]string{
		types.ModuleName: nil,
	}
	sk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)
	govK := government.NewKeeper(cdc, keys[government.StoreKey])

	// Set the government address
	_ = govK.SetGovernmentAddress(ctx, TestGovernment)

	// Setup the Did Document
	TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	TestDidDocument = setupDidDocument(ctx, ak, "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

	idk := NewKeeper(cdc, keys[types.StoreKey], ak, sk)

	// Set initial supply
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 1))))

	// Set module accounts
	idAcc := supply.NewEmptyModuleAccount(types.ModuleName)
	bech32Addr, err := bech32.ConvertAndEncode("did:com:", idAcc.Address.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Print(bech32Addr)

	idk.supplyKeeper.SetModuleAccount(ctx, idAcc)

	return cdc, ctx, ak, bk, govK, idk
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	government.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
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
		Context: "https://www.w3.org/ns/did/v1",
		ID:      address,
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
				ID:           fmt.Sprintf("%s#keys-1", address),
				Type:         "Secp256k1VerificationKey2018",
				Controller:   address,
				PublicKeyHex: hex.EncodeToString(secp256k1Key[:]),
			},
			types.PubKey{
				ID:           fmt.Sprintf("%s#keys-2", address),
				Type:         "RsaVerificationKey2018",
				Controller:   address,
				PublicKeyHex: "04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a",
			},
		},
	}
}

// Identities
var TestOwnerAddress sdk.AccAddress
var TestDidDocument types.DidDocument

// Deposit requests
var TestGovernment, _ = sdk.AccAddressFromBech32("cosmos1gdpsu89prllyw49eehskv6t8800p6chefyuuwe")
var TestDepositor, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var TestPairwiseDid, _ = sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")
var TestDidDepositRequest = types.DidDepositRequest{
	Recipient:     TestPairwiseDid,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	FromAddress:   TestDepositor,
}

// Power up requests
var TestDidPowerUpRequest = types.DidPowerUpRequest{
	Claimant:      TestDepositor,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
}
