package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	docsStoreKey sdk.StoreKey
	cdc          *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		docsStoreKey: storeKey,
		cdc:          cdc,
	}
}

// ----------------------------------
// --- Keeper methods
// ----------------------------------

// GetMetadata returns the Metadata Reference for the document with the given Reference.
func (keeper Keeper) GetSharedDocumentsWithUser(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress) string {
	//store := ctx.KVStore(keeper.docsStoreKey)
	//result := store.Get([]byte(reference))
	//return string(result)
	return ""
}

// ShareDocument allows the sharing of a document represented by the given Reference, between the given sender and the
// given recipient.
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) sdk.Error {
	//store := ctx.KVStore(keeper.docsStoreKey)
	return nil
}
