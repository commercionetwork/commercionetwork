package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/commercionetwork/commercionetwork/x/did/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func getTimestamp(ctx sdk.Context) (*time.Time, error) {
	// Following the W3C Decentralized Identifiers (DIDs) v1.0 for
	// XML Datetime normalized to UTC 00:00:00 and without sub-second decimal precision
	// format := "2020-12-20T19:17:47Z"

	t := ctx.BlockTime()

	// t, err := time.Parse(format, tString)
	// if err != nil {
	// 	return nil, err
	// }
	return &t, nil
}

// AppendDid appends a DID document in the store with given id
func (k Keeper) AppendDid(ctx sdk.Context, didDocument types.DidDocument) string {
	// Create the Document
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(sdk.AccAddress(didDocument.ID)), k.cdc.MustMarshalBinaryBare(&didDocument))
	return didDocument.ID
}

// GetDdoByOwner returns the DID document reference associated to a given DID.
// If the given DID has no DID document reference associated, returns nil.
func (k Keeper) GetDdoByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, fmt.Errorf("DID document with owner %s not found", owner.String())
	}

	var DidDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &DidDocument)
	return DidDocument, nil
}

func getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

func (k Keeper) HasIdentity(ctx sdk.Context, ID string) bool {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(sdk.AccAddress(ID))
	return store.Has(identityKey)
}

// GetAllDidDocument returns all DID document
func (k Keeper) GetAllDidDocument(ctx sdk.Context) (list []types.DidDocument) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.IdentitiesStorePrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DidDocument
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
