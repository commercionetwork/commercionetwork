package commercioauth

import (
	app "commercio-network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"testing"
)

var keyValue = "A8dJWr6t9Yh31YYvXkb0N/HtkC5J+KAP75dqg8pr3uws"
var address = "cosmos153eu7p9lpgaatml7ua2vvgl8w08r4kjl5ca3y0"
var keyType = "Secp256k1"

var k = Keeper{
	accountKeeper: auth.AccountKeeper{},
	cdc:           testCdc,
}

func TestKeeper_RegisterAccount(t *testing.T) {
	actual := k.RegisterAccount(ctx, address, keyType, keyValue)

	if actual != nil {
		t.Errorf("Registration of account %s failed", actual)
	}

}
