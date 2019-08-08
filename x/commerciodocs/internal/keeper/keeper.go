package keeper

import (
	"bytes"
	"fmt"
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/utilities"
	"github.com/commercionetwork/commercionetwork/x/commercioid"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/genproto/googleapis/type/date"
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

func (keeper Keeper) GetDocument(ctx sdk.Context, checksumValue string) types.Document {

}

func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, sender sdk.AccAddress) []types.Document {

}

func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, sender sdk.AccAddress) []types.Document {

}

// GetSharedDocumentsWithUser
func (keeper Keeper) GetSharedDocumentsWithUser(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress) []types.Document {

}

// ShareDocument allows the sharing of a document represented by the given Reference, between the given sender and the
// given recipient.
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) sdk.Error {

	store := ctx.KVStore(keeper.docsStoreKey)

	return nil
}
