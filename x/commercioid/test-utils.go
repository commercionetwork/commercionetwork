package commercioid

/*
import (
	"commercio-network/types"
	"commercio-network/x/commercioid/internal/keeper"
	types2 "commercio-network/x/commercioid/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	cdc        *codec.Codec
	ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.BaseKeeper
	idKeeper   keeper.Keeper
}

//commercioauth module initialisation
var input = setupTestInput()

//This function create an enviroment to test modules
func setupTestInput() testInput {

	db := dbm.NewMemDB()
	cdc := makeCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// CommercioID
	keyIDIdentities := sdk.NewKVStoreKey("id_identities")
	keyIDOwners := sdk.NewKVStoreKey("id_owners")
	keyIDConnections := sdk.NewKVStoreKey("id_connections")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyIDIdentities, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDOwners, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyIDConnections, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	idk := keeper.NewKeeper(keyIDIdentities, keyIDOwners, keyIDConnections, cdc)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{
		cdc:        cdc,
		ctx:        ctx,
		accKeeper:  ak,
		bankKeeper: bk,
		idKeeper:   idk,
	}

}

func makeCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(types2.MsgSetIdentity{}, "commercioid/SetIdentity", nil)
	cdc.RegisterConcrete(types2.MsgCreateConnection{}, "commectioid/CreateConnection", nil)

	cdc.Seal()

	return cdc
}

//TEST VARS

var msgSetId = types2.MsgSetIdentity{
	DID:          ownerIdentity,
	DDOReference: identityRef,
	Owner:        owner,
}

var msgCreateConn = types2.MsgCreateConnection{
	FirstUser:  ownerIdentity,
	SecondUser: recipient,
	Signer:     owner,
}

var address = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var owner, _ = sdk.AccAddressFromBech32(address)
var ownerIdentity = types.Did("newReader")
var identityRef = "ddo-reference"
var reference = "reference"
var metadata = "metadata"
var recipient = types.Did("recipient")

*/
