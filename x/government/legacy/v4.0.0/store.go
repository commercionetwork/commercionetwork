package v400

import (
	v300 "github.com/commercionetwork/commercionetwork/x/government/legacy/v3.0.0"
	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {

	//newCTX := sdk.NewContext()

	oldKVStoreKey := sdk.NewKVStoreKey(v300.ModuleName)
	ctx.MultiStore()
	oldstore := ctx.KVStore(oldKVStoreKey)
	store := ctx.KVStore(storeKey)
	v300key := []byte(v300.GovernmentStoreKey)

	if !oldstore.Has(v300key) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Gov address not present")
	}
	migrateGovKeys(store, oldstore)
	return nil
}

func migrateGovKeys(store sdk.KVStore, oldstore sdk.KVStore) {
	v300key := []byte(v300.GovernmentStoreKey)
	govValue := oldstore.Get(v300key)
	store.Set([]byte(types.GovernmentStoreKey), govValue)
	oldstore.Delete(v300key)
}
