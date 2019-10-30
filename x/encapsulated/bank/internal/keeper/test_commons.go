package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

// SetupTestInput creates a test environment
func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()

	key := sdk.NewKVStoreKey("bank")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())
	k := NewKeeper(cdc, key, nil)

	return cdc, ctx, k
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)
	cdc.Seal()

	return cdc
}
