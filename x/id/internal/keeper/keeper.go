package keeper

import (
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	StoreKey   sdk.StoreKey
	cdc        *codec.Codec
	bankKeeper bank.Keeper
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, bankKeeper bank.Keeper) Keeper {
	return Keeper{
		StoreKey:   storeKey,
		cdc:        cdc,
		bankKeeper: bankKeeper,
	}
}

// ------------------
// --- Identities
// ------------------

func (keeper Keeper) getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return []byte(types.IdentitiesStorePrefix + owner.String())
}

// SaveIdentity saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (keeper Keeper) SaveIdentity(ctx sdk.Context, owner sdk.AccAddress, document types.DidDocument) {
	identitiesStore := ctx.KVStore(keeper.StoreKey)
	identitiesStore.Set(keeper.getIdentityStoreKey(owner), keeper.cdc.MustMarshalBinaryBare(document))
}

// GetDidDocumentByOwner returns the Did Document reference associated to a given Did
func (keeper Keeper) GetDidDocumentByOwner(ctx sdk.Context, owner sdk.AccAddress) (didDocument types.DidDocument, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	identityKey := keeper.getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, true
}

// GetIdentities returns the list of all identities for the given context
func (keeper Keeper) GetIdentities(ctx sdk.Context) ([]types.Identity, error) {
	store := ctx.KVStore(keeper.StoreKey)
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

// ----------------------------
// --- Did deposit requests
// ----------------------------

func (keeper Keeper) getDepositRequestStoreKey(proof string) []byte {
	return []byte(types.DidDepositRequestStorePrefix + proof)
}

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (keeper Keeper) StoreDidDepositRequest(ctx sdk.Context, request types.DidDepositRequest) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	requestKey := keeper.getDepositRequestStoreKey(request.Proof)
	if store.Has(requestKey) {
		return sdk.ErrUnknownRequest("Did deposit request with the same proof already exists")
	}

	store.Set(requestKey, keeper.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetDidDepositRequestByProof returns the request having the same proof.
func (keeper Keeper) GetDidDepositRequestByProof(ctx sdk.Context, proof string) (request types.DidDepositRequest, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	requestKey := keeper.getDepositRequestStoreKey(proof)
	if !store.Has(requestKey) {
		return types.DidDepositRequest{}, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(store.Get(requestKey), &request)
	return request, true
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (keeper Keeper) ChangeDepositRequestStatus(ctx sdk.Context, proof string, status types.RequestStatus) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	request, found := keeper.GetDidDepositRequestByProof(ctx, proof)
	if !found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Did deposit request with proof %s not fond", proof))
	}

	// Update and store the request
	request.Status = &status
	store.Set(keeper.getDepositRequestStoreKey(request.Proof), keeper.cdc.MustMarshalBinaryBare(request))

	return nil
}

// GetDepositRequests returns the list of the deposit requests existing inside the given context
func (keeper Keeper) GetDepositRequests(ctx sdk.Context) (requests []types.DidDepositRequest) {
	store := ctx.KVStore(keeper.StoreKey)

	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidDepositRequestStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidDepositRequest
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}

// ----------------------------
// --- Did power up requests
// ----------------------------

func (keeper Keeper) getDidPowerUpRequestStoreKey(proof string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + proof)
}

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (keeper Keeper) StorePowerUpRequest(ctx sdk.Context, request types.DidPowerUpRequest) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	requestStoreKey := keeper.getDidPowerUpRequestStoreKey(request.Proof)
	if store.Has(requestStoreKey) {
		return sdk.ErrUnknownRequest("PowerUp request with the same proof already exists")
	}

	store.Set(requestStoreKey, keeper.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetDidDepositRequestByProof returns the request having the same proof.
func (keeper Keeper) GetPowerUpRequestByProof(ctx sdk.Context, proof string) (request types.DidPowerUpRequest, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	requestStoreKey := keeper.getDidPowerUpRequestStoreKey(proof)
	if !store.Has(requestStoreKey) {
		return types.DidPowerUpRequest{}, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(store.Get(requestStoreKey), &request)
	return request, true
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (keeper Keeper) ChangePowerUpRequestStatus(ctx sdk.Context, proof string, status types.RequestStatus) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	request, found := keeper.GetPowerUpRequestByProof(ctx, proof)
	if !found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("PowerUp request with proof %s not found", proof))
	}

	// Update and store the request
	request.Status = &status
	store.Set(keeper.getDidPowerUpRequestStoreKey(proof), keeper.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetPowerUpRequests returns the list the requests saved inside the given context
func (keeper Keeper) GetPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	store := ctx.KVStore(keeper.StoreKey)

	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidPowerUpRequestStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidPowerUpRequest
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}

// ------------------------
// --- Deposits handling
// ------------------------

// DepositIntoPool allows to deposit the specified amount into the liquidity pool, taking it from the
// specified depositor balance
func (keeper Keeper) DepositIntoPool(ctx sdk.Context, depositor sdk.AccAddress, amount sdk.Coins) sdk.Error {
	// Check the amount
	if !amount.IsValid() || amount.Empty() || amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(amount.String())
	}

	// Subtract the coins from the user
	if _, err := keeper.bankKeeper.SubtractCoins(ctx, depositor, amount); err != nil {
		return err
	}

	// Get the current pool
	currentPool := keeper.GetPoolAmount(ctx)
	if err := keeper.SetPoolAmount(ctx, currentPool.Add(amount)); err != nil {
		return err
	}

	return nil
}

// FundAccount allows to take the specified amount from the liquidity pool and move them into the
// specified account balance
func (keeper Keeper) FundAccount(ctx sdk.Context, account sdk.AccAddress, amount sdk.Coins) sdk.Error {
	// Check the amount
	if !amount.IsValid() || amount.Empty() || amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(amount.String())
	}

	// Get the current pool
	currentPool := keeper.GetPoolAmount(ctx)

	// Check that the pool has enough funds
	if amount.IsAnyGT(currentPool) {
		return sdk.ErrInsufficientFunds("Pool does not have enough funds")
	}

	// Update the pool
	currentPool = currentPool.Sub(amount)
	if err := keeper.SetPoolAmount(ctx, currentPool); err != nil {
		return err
	}

	// Add the coins to the user
	if _, err := keeper.bankKeeper.AddCoins(ctx, account, amount); err != nil {
		return err
	}

	return nil
}

// SetPoolAmount allows to set the pool amount to the given one
func (keeper Keeper) SetPoolAmount(ctx sdk.Context, amount sdk.Coins) sdk.Error {
	if amount == nil || !amount.IsValid() || amount.IsAnyNegative() {
		return sdk.ErrInvalidCoins("Invalid pool amount")
	}

	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.DepositsPoolStoreKey), keeper.cdc.MustMarshalBinaryBare(&amount))

	return nil
}

// GetPoolAmount returns the current pool amount
func (keeper Keeper) GetPoolAmount(ctx sdk.Context) (pool sdk.Coins) {
	store := ctx.KVStore(keeper.StoreKey)
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.DepositsPoolStoreKey)), &pool)
	return pool
}
