package commercioauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

var keyValue = "A8dJWr6t9Yh31YYvXkb0N/HtkC5J+KAP75dqg8pr3uws"
var address = "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0"
var keyType = "Secp256k1"

var input = setupTestInput()
var keeper = Keeper{
	accountKeeper: input.accKeeper,
	cdc:           input.cdc,
}

//All fields are valid
func TestKeeper_RegisterAccount(t *testing.T) {

	ass := assert.New(t)

	actual := keeper.RegisterAccount(input.ctx, address, keyType, keyValue)

	if !ass.Nil(actual) {
		t.Errorf("Registration of account %s failed", actual)
	}
}

//Invalid address
func TestKeeper_RegisterAccount2(t *testing.T) {

	ass := assert.New(t)

	invalidAddr := "cosmos153eu7p9lpgaetml7ua2vvgl8w08r4kjl5ca3y2"
	actual := keeper.RegisterAccount(input.ctx, invalidAddr, keyType, keyValue)

	expected := sdk.ErrInvalidAddress("Invalid address provided")

	if !ass.Equal(actual.Error(), expected.Error()) {
		t.Errorf("account address should be invalid")
	}
}

//Invalid key type
func TestKeeper_RegisterAccount3(t *testing.T) {

	ass := assert.New(t)

	invalidKeyType := "testKey"
	key := "64696f"
	actual := keeper.RegisterAccount(input.ctx, address, invalidKeyType, key)

	expected := sdk.ErrUnknownRequest("Invalid key type. Currently supported key types are Ed25519 and Secp256k1")

	if !ass.Equal(expected.Error(), actual.Error()) {
		t.Errorf("key type should be invalid")
	}
}

//Invalid key value
func TestKeeper_RegisterAccount4(t *testing.T) {

	ass := assert.New(t)

	invalidKeyValue := "key"
	actual := keeper.RegisterAccount(input.ctx, address, keyType, invalidKeyValue)

	expected := sdk.ErrInternal("Can't set a null public key to account cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0")

	if !ass.Equal(expected, actual) {
		t.Errorf(actual.Error())
	}
}

//Invalid address
func TestKeeper_GetAccount(t *testing.T) {
	ass := assert.New(t)

	invalidAddr := "cosmos153eu7p9lpgaetml7ua2vvgl8w08r4kjl5ca3y2"
	expected := sdk.ErrInvalidAddress("Invalid address provided")

	_, actualErr := keeper.GetAccount(input.ctx, invalidAddr)

	if !ass.Equal(expected, actualErr) {
		t.Errorf("The two errors should be equals because address is invalid")
	}
}

//Account not found
func TestKeeper_GetAccount2(t *testing.T) {
	ass := assert.New(t)

	expected := sdk.ErrInvalidAddress("No account found for address cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0")

	_, actualErr := keeper.GetAccount(input.ctx, address)

	if !ass.Equal(expected.Error(), actualErr.Error()) {
		t.Errorf("The two errors should be equals because address isnt found")
	}
}

//Account found
func TestKeeper_GetAccount3(t *testing.T) {
	ass := assert.New(t)

	addr, err := sdk.AccAddressFromBech32(address)

	if err != nil {
		panic(err)
	}

	account := keeper.accountKeeper.NewAccountWithAddress(input.ctx, addr)
	account.SetPubKey(ed25519.GenPrivKey().PubKey())
	keeper.accountKeeper.SetAccount(input.ctx, account)

	actual, actualErr := keeper.GetAccount(input.ctx, address)

	if actualErr != nil {
		panic(actualErr)
	}

	if ass.NotEqual(account.String(), actual.String()) {
		t.Errorf("The two account should be equals")
	}
}

func TestKeeper_ListAccounts(t *testing.T) {
	ass := assert.New(t)

	var expected []auth.Account

	actual := keeper.ListAccounts(input.ctx)

	if !ass.Equal(expected, actual) {
		t.Errorf("The two slices are empty and should be equals")
	}
}
