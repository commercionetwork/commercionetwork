package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Keeper struct {
	storeKey  sdk.StoreKey
	cdc       *codec.Codec
	accKeeper auth.AccountKeeper
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accKeeper auth.AccountKeeper) Keeper {
	return Keeper{
		storeKey:  storeKey,
		cdc:       cdc,
		accKeeper: accKeeper,
	}
}

func (keeper Keeper) getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return []byte(types.IdentitiesStorePrefix + owner.String())
}

// SaveIdentity saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (keeper Keeper) SaveIdentity(ctx sdk.Context, owner sdk.AccAddress, document types.DidDocument) sdk.Error {
	account := keeper.accKeeper.GetAccount(ctx, owner)
	accountPubKey := account.GetPubKey()

	// --------------------------------------------
	// --- Check the authentication key validity
	// --------------------------------------------

	// Get the authentication key
	authKey, found := document.PubKeys.FindById(document.Authentication[0])
	if !found {
		return sdk.ErrUnknownRequest("Authentication key not found inside publicKey array")
	}

	if _, isEd25519 := accountPubKey.(ed25519.PubKeyEd25519); isEd25519 {
		// Check the auth key type coherence with its value
		if authKey.Type != types.KeyTypeEd25519 {
			msg := fmt.Sprintf("Invalid authentication key value, must be of type %s", types.KeyTypeEd25519)
			return sdk.ErrUnknownRequest(msg)
		}
	}

	if _, isSecp256k1 := accountPubKey.(secp256k1.PubKeySecp256k1); isSecp256k1 {
		// Check the auth key type coherence with its value
		if authKey.Type != types.KeyTypeSecp256k1 {
			msg := fmt.Sprintf("Invalid authentication key value, must be of type %s", types.KeyTypeSecp256k1)
			return sdk.ErrUnknownRequest(msg)
		}
	}

	// Get the authentication key bytes
	authKeyBytes, err := hex.DecodeString(authKey.PublicKeyHex)
	if err != nil {
		return sdk.ErrUnknownRequest("Invalid authentication key hex value")
	}

	// Check that the authentication key bytes are the same of the key associated with the account
	if !bytes.Equal(accountPubKey.Bytes()[:], authKeyBytes[:]) {
		return sdk.ErrUnknownRequest("Authentication key is not the one associated with the account")
	}

	// TODO: Check that the proof signatureValue is the valid signature of the entire Did Document made with the user private key

	// ------------------------------
	// --- Store the Did Document
	// ------------------------------

	// Set the Did Document into the store
	identitiesStore := ctx.KVStore(keeper.storeKey)
	identitiesStore.Set(keeper.getIdentityStoreKey(owner), keeper.cdc.MustMarshalBinaryBare(document))

	return nil
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (keeper Keeper) GetDidDocumentByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, bool) {
	store := ctx.KVStore(keeper.storeKey)

	identityKey := keeper.getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, false
	}

	var didDocument types.DidDocument
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, true
}

// -------------------------
// --- Genesis utils
// -------------------------

// GetIdentities returns the list of all identities for the given context
func (keeper Keeper) GetIdentities(ctx sdk.Context) ([]types.Identity, error) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	var identities []types.Identity
	for ; iterator.Valid(); iterator.Next() {
		addressString := strings.ReplaceAll(string(iterator.Key()), types.IdentitiesStorePrefix, "")
		address, err := sdk.AccAddressFromBech32(addressString)
		if err != nil {
			return nil, err
		}

		var didDocument types.DidDocument
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &didDocument)

		identity := types.Identity{
			Owner:       address,
			DidDocument: didDocument,
		}
		identities = append(identities, identity)
	}

	return identities, nil
}

// SetIdentities allows to bulk save a bunch of identities inside the store
func (keeper Keeper) SetIdentities(ctx sdk.Context, identities []types.Identity) {
	for _, identity := range identities {
		keeper.SaveIdentity(ctx, identity.Owner, identity.DidDocument)
	}
}
