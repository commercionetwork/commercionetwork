package commercioauth

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	accountKeeper auth.AccountKeeper

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	cdc *codec.Codec
}

// NewKeeper creates new instances of the CommercioAUTH Keeper
func NewKeeper(
	accountKeeper auth.AccountKeeper,
	cdc *codec.Codec) Keeper {
	return Keeper{
		accountKeeper: accountKeeper,
		cdc:           cdc,
	}
}

// ----------------------------------
// --- Keeper methods
// ----------------------------------

const KeyTypeEd25519 = "Ed25519"
const KeyTypeEcdsa = "Secp256k1"

func (keeper Keeper) RegisterAccount(ctx sdk.Context, address string, keyType string, keyValue string) sdk.Error {

	// Get the address from the string
	accountAddress, err := sdk.AccAddressFromHex(address)
	if err != nil {
		return sdk.ErrInvalidAddress("Invalid address provided")
	}

	// Decode the HEX string
	publicKeyBytes, err := hex.DecodeString(keyValue)
	if err != nil {
		return sdk.ErrUnknownRequest("Invalid public key HEX value provided")
	}

	// Validate the sent key type
	var publicKey crypto.PubKey
	if keyType == KeyTypeEd25519 {
		var pkEd ed25519.PubKeyEd25519
		copy(pkEd[:], publicKeyBytes[:])
		publicKey = pkEd
	} else if keyType == KeyTypeEcdsa {
		var pkEd secp256k1.PubKeySecp256k1
		copy(pkEd[:], publicKeyBytes[:])
		publicKey = pkEd
	} else {
		return sdk.ErrUnknownRequest("Invalid key type. Currently supported key types are Ed25519 and Secp256k1")
	}

	// Try getting the existing account
	account := keeper.accountKeeper.GetAccount(ctx, accountAddress)

	if account == nil {
		// Create a new account
		account = keeper.accountKeeper.NewAccountWithAddress(ctx, accountAddress)

		// Set the account's public key
		if err := account.SetPubKey(publicKey); err != nil {
			return sdk.ErrInternal(fmt.Sprintf("Error while setting the account's public key %s", err))
		}

		// Store the account inside the store
		keeper.accountKeeper.SetAccount(ctx, account)
	}

	return nil
}

func (keeper Keeper) GetAccount(ctx sdk.Context, address string) (auth.Account, sdk.Error) {

	// Get the address from the string
	accountAddress, err := sdk.AccAddressFromHex(address)
	if err != nil {
		return nil, sdk.ErrInvalidAddress("Invalid address provided")
	}

	// Try getting the existing account
	var accountError sdk.Error = nil
	account := keeper.accountKeeper.GetAccount(ctx, accountAddress)
	if account == nil {
		accountError = sdk.ErrInvalidAddress(fmt.Sprintf("No account found for address %s", address))
	}
	return account, accountError
}

func (keeper Keeper) ListAccounts(ctx sdk.Context) []auth.Account {
	var accounts []auth.Account
	appendAccount := func(acc auth.Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	keeper.accountKeeper.IterateAccounts(ctx, appendAccount)

	return accounts
}
