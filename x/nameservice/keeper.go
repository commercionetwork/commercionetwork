package nameservice

import (
	// Provides tools to work with the Cosmos encoding format, Amino.
	"github.com/cosmos/cosmos-sdk/codec"

	// Controls accounts and coin transfers.
	"github.com/cosmos/cosmos-sdk/x/bank"

	// Contains commonly used types throughout the SDK.
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {

	// Reference to the Keeper from the bank module. Including it allows code in this module to call functions from the
	// bank module. The SDK uses an object capabilities approach to accessing sections of the application state.
	// This is to allow developers to employ a least authority approach, limiting the capabilities of a faulty or
	// malicious module from affecting parts of state it doesn't need access to.
	coinKeeper bank.Keeper

	namesStoreKey  sdk.StoreKey // Unexposed key to access name store from sdk.Context
	ownersStoreKey sdk.StoreKey // Unexposed key to access owners store from sdk.Context
	pricesStoreKey sdk.StoreKey // Unexposed key to access prices store from sdk.Context

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	cdc *codec.Codec
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	// Get the store object for the map[name] value using the the namesStoreKey from the Keeper
	store := ctx.KVStore(k.namesStoreKey)

	// Insert the <name, value> pair into the store using its .Set([]byte, []byte) method.
	// As the store only takes []byte, first cast the strings to []byte and the use them as parameters into the Set method.
	store.Set([]byte(name), []byte(value))
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	// Access the store using the StoreKey
	store := ctx.KVStore(k.namesStoreKey)

	// Use the .Get([]byte) []byte method. As the parameter into the function, pass the key, which is the name string
	// casted to []byte, and get back the result in the form of []byte
	bz := store.Get([]byte(name))

	// Cast this to a string and return the result.
	return string(bz)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.ownersStoreKey)
	bz := store.Get([]byte(name))
	return bz != nil
}

// GetOwner - get the current owner of a name
func (k Keeper) GetOwner(ctx sdk.Context, name string) sdk.AccAddress {
	store := ctx.KVStore(k.ownersStoreKey)
	bz := store.Get([]byte(name))
	return bz
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	store := ctx.KVStore(k.ownersStoreKey)

	// Because sdk.AccAddress is a type alias for []byte, it can natively be casted to it
	store.Set([]byte(name), owner)
}

// GetPrice - gets the current price of a name.  If price doesn't exist yet, set to 1steak.
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	if !k.HasOwner(ctx, name) {
		return sdk.Coins{sdk.NewInt64Coin("mycoin", 1)}
	}

	store := ctx.KVStore(k.pricesStoreKey)
	bz := store.Get([]byte(name))
	var price sdk.Coins
	k.cdc.MustUnmarshalBinaryBare(bz, &price)
	return price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	store := ctx.KVStore(k.pricesStoreKey)
	store.Set([]byte(name), k.cdc.MustMarshalBinaryBare(price))
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(
	coinKeeper bank.Keeper,
	namesStoreKey sdk.StoreKey,
	ownersStoreKey sdk.StoreKey,
	priceStoreKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper:     coinKeeper,
		namesStoreKey:  namesStoreKey,
		ownersStoreKey: ownersStoreKey,
		pricesStoreKey: priceStoreKey,
		cdc:            cdc,
	}
}
