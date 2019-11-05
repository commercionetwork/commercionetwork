package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	accKeeper    auth.AccountKeeper
	supplyKeeper supply.Keeper
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accKeeper auth.AccountKeeper, supplyKeeper supply.Keeper) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		accKeeper:    accKeeper,
		supplyKeeper: supplyKeeper,
	}
}

// ------------------
// --- Identities
// ------------------

func (k Keeper) getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return []byte(types.IdentitiesStorePrefix + owner.String())
}

// SaveDidDocument saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (k Keeper) SaveDidDocument(ctx sdk.Context, document types.DidDocument) sdk.Error {
	owner := document.ID

	// Get the account and its public key
	account := k.accKeeper.GetAccount(ctx, owner)
	if account == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Could not find account %s", owner.String()))
	}

	accountPubKey := account.GetPubKey()

	// --------------------------------------------
	// --- Check the authentication key validity
	// --------------------------------------------

	// Get the authentication key
	authKey, found := document.PubKeys.FindByID(document.Authentication[0])
	if !found {
		return sdk.ErrUnknownRequest("Authentication key not found inside publicKey array")
	}

	// Get the authentication key bytes
	authKeyBytes, err := hex.DecodeString(authKey.PublicKeyHex)
	if err != nil {
		return sdk.ErrUnknownRequest("Invalid authentication key hex value")
	}

	// TODO: The following code should be tested

	var key []byte
	if ed25519key, isEd25519 := accountPubKey.(ed25519.PubKeyEd25519); isEd25519 {
		// Check the auth key type coherence with its value
		if authKey.Type != types.KeyTypeEd25519 {
			msg := fmt.Sprintf("Invalid authentication key value, must be of type %s", types.KeyTypeEd25519)
			return sdk.ErrUnknownRequest(msg)
		}
		key = ed25519key[:]
	}

	if secp256k1key, isSecp256k1 := accountPubKey.(secp256k1.PubKeySecp256k1); isSecp256k1 {
		// Check the auth key type coherence with its value
		if authKey.Type != types.KeyTypeSecp256k1 {
			msg := fmt.Sprintf("Invalid authentication key value, must be of type %s", types.KeyTypeSecp256k1)
			return sdk.ErrUnknownRequest(msg)
		}
		key = secp256k1key[:]
	}

	// Check that the authentication key bytes are the same of the key associated with the account
	if !bytes.Equal(authKeyBytes, key) {
		msg := fmt.Sprintf(
			"Authentication key is not the one associated with the account. Expected %s but got %s",
			hex.EncodeToString(key),
			hex.EncodeToString(authKeyBytes),
		)
		return sdk.ErrUnknownRequest(msg)
	}

	// TODO: Check that the proof signatureValue is the valid signature of the entire Did Document made with the user private key

	// ------------------------------
	// --- Store the Did Document
	// ------------------------------

	// Set the Did Document into the store
	identitiesStore := ctx.KVStore(k.storeKey)
	identitiesStore.Set(k.getIdentityStoreKey(owner), k.cdc.MustMarshalBinaryBare(document))

	return nil
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (k Keeper) GetDidDocumentByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, bool) {
	store := ctx.KVStore(k.storeKey)

	identityKey := k.getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, false
	}

	var didDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, true
}

// -------------------------
// --- Genesis utils
// -------------------------

// GetDidDocuments returns the list of all identities for the given context
func (k Keeper) GetDidDocuments(ctx sdk.Context) ([]types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	var didDocuments []types.DidDocument
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var didDocument types.DidDocument
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &didDocument)
		didDocuments = append(didDocuments, didDocument)
	}

	return didDocuments, nil
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

func (k Keeper) getDepositRequestStoreKey(proof string) []byte {
	return []byte(types.DidDepositRequestStorePrefix + proof)
}

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (k Keeper) StoreDidDepositRequest(ctx sdk.Context, request types.DidDepositRequest) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	requestKey := k.getDepositRequestStoreKey(request.Proof)
	if store.Has(requestKey) {
		return sdk.ErrUnknownRequest("Did deposit request with the same proof already exists")
	}

	store.Set(requestKey, k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetDidDepositRequestByProof returns the request having the same proof.
func (k Keeper) GetDidDepositRequestByProof(ctx sdk.Context, proof string) (request types.DidDepositRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	requestKey := k.getDepositRequestStoreKey(proof)
	if !store.Has(requestKey) {
		return types.DidDepositRequest{}, false
	}

	k.cdc.MustUnmarshalBinaryBare(store.Get(requestKey), &request)
	return request, true
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (k Keeper) ChangeDepositRequestStatus(ctx sdk.Context, proof string, status types.RequestStatus) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	request, found := k.GetDidDepositRequestByProof(ctx, proof)
	if !found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Did deposit request with proof %s not fond", proof))
	}

	// Update and store the request
	request.Status = &status
	store.Set(k.getDepositRequestStoreKey(request.Proof), k.cdc.MustMarshalBinaryBare(request))

	return nil
}

// GetDepositRequests returns the list of the deposit requests existing inside the given context
func (k Keeper) GetDepositRequests(ctx sdk.Context) (requests []types.DidDepositRequest) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidDepositRequestStorePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidDepositRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}

// ----------------------------
// --- Did power up requests
// ----------------------------

func (k Keeper) getDidPowerUpRequestStoreKey(proof string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + proof)
}

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (k Keeper) StorePowerUpRequest(ctx sdk.Context, request types.DidPowerUpRequest) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	requestStoreKey := k.getDidPowerUpRequestStoreKey(request.Proof)
	if store.Has(requestStoreKey) {
		return sdk.ErrUnknownRequest("PowerUp request with the same proof already exists")
	}

	store.Set(requestStoreKey, k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetDidDepositRequestByProof returns the request having the same proof.
func (k Keeper) GetPowerUpRequestByProof(ctx sdk.Context, proof string) (request types.DidPowerUpRequest, found bool) {
	store := ctx.KVStore(k.storeKey)

	requestStoreKey := k.getDidPowerUpRequestStoreKey(proof)
	if !store.Has(requestStoreKey) {
		return types.DidPowerUpRequest{}, false
	}

	k.cdc.MustUnmarshalBinaryBare(store.Get(requestStoreKey), &request)
	return request, true
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (k Keeper) ChangePowerUpRequestStatus(ctx sdk.Context, proof string, status types.RequestStatus) sdk.Error {
	store := ctx.KVStore(k.storeKey)

	request, found := k.GetPowerUpRequestByProof(ctx, proof)
	if !found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("PowerUp request with proof %s not found", proof))
	}

	// Update and store the request
	request.Status = &status
	store.Set(k.getDidPowerUpRequestStoreKey(proof), k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetPowerUpRequests returns the list the requests saved inside the given context
func (k Keeper) GetPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidPowerUpRequestStorePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidPowerUpRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}

// ------------------------
// --- Deposits handling
// ------------------------

func (k Keeper) SetPowerUpRequestHandled(ctx sdk.Context, activationReference string) {
	handledRequests := k.GetHandledPowerUpRequestsReferences(ctx)
	if requests, edited := handledRequests.AppendIfMissing(activationReference); edited {
		k.SetHandledPowerUpRequestsReferences(ctx, requests)
	}
}

func (k Keeper) GetHandledPowerUpRequestsReferences(ctx sdk.Context) ctypes.Strings {
	store := ctx.KVStore(k.storeKey)

	var handledRequests ctypes.Strings
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.HandledPowerUpRequestsStoreKey)), &handledRequests)
	return handledRequests
}

func (k Keeper) SetHandledPowerUpRequestsReferences(ctx sdk.Context, references ctypes.Strings) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.HandledPowerUpRequestsStoreKey), k.cdc.MustMarshalBinaryBare(&references))
}

// DepositIntoPool allows to deposit the specified amount into the liquidity pool, taking it from the
// specified depositor balance
func (k Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) sdk.Error {
	// Check the amount
	if !amount.IsValid() || amount.Empty() || amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("Invalid coins: %s", amount))
	}

	// Subtract the coins from the user
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositor, types.ModuleName, amount); err != nil {
		return err
	}

	return nil
}

// FundAccount allows to take the specified amount from the liquidity pool and move them into the
// specified account balance
func (k Keeper) FundAccount(ctx sdk.Context, account sdk.AccAddress, amount sdk.Coins) sdk.Error {
	// Check the amount
	if amount.Empty() || !amount.IsValid() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("Invalid coins: %s", amount))
	}

	// Get the current pool
	currentPool := k.GetPoolAmount(ctx)

	// Check that the pool has enough funds
	if amount.IsAnyGT(currentPool) {
		return sdk.ErrInsufficientFunds("Pool does not have enough funds")
	}

	// Add the coins to the user
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, account, amount); err != nil {
		return err
	}

	return nil
}
