package v400

import (
	v300 "github.com/commercionetwork/commercionetwork/x/government/legacy/v3.0.0"
	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	v300key := []byte(v300.GovernmentStoreKey)

	if !store.Has(v300key) {
		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "Gov address not present")
	}
	migrateGovKeys(store)
	return nil
}

func migrateGovKeys(store storetypes.KVStore) {
	v300key := []byte(v300.GovernmentStoreKey)
	govValue := store.Get(v300key)
	store.Set([]byte(types.GovernmentStoreKey), govValue)
	store.Delete(v300key)
}
