package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func init() {
	configTestPrefixes()
	governmentTestAddress, _ = sdk.AccAddressFromBech32("did:com:1zg4jreq2g57s4efrl7wnh2swtrz3jt9nfaumcm")
	notGovernmentAddress, _ = sdk.AccAddressFromBech32("did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd")
}

func configTestPrefixes() {
	AccountAddressPrefix := "did:com:"
	AccountPubKeyPrefix := AccountAddressPrefix + "pub"
	ValidatorAddressPrefix := AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix := AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

var governmentTestAddress, notGovernmentAddress sdk.AccAddress

// This function creates an environment to test the government module
// if address is defined it will be used to add the government address
func setupKeeperWithGovernmentAddress(t testing.TB, address sdk.AccAddress) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(
		codec.NewProtoCodec(registry), storeKey, memStoreKey,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	if address != nil {
		store := ctx.KVStore(keeper.storeKey)
		store.Set([]byte(types.GovernmentStoreKey), address)
	}

	return keeper, ctx
}
