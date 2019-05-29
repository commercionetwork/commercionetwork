package commercioauth

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
	"testing"

	id "commercio-network/x/commercioid"
)

var comAuthkeeper = Keeper{
	accountKeeper: input.accKeeper,
	cdc:           input.cdc,
}

//handled message
func TestNewHandler(t *testing.T) {

	ass := assert.New(t)

	handler := NewHandler(comAuthkeeper)

	msg := MsgCreateAccount{
		Signer:   sdk.AccAddress(priv.PubKey().Address()),
		Address:  "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0",
		KeyType:  "Secp256k1",
		KeyValue: "cosmospub1addwnpepqvhd7lalnwuh0k2hy8x60k99dhqurv02eeuv3cencr2c088gp03uq6wrh5f",
	}

	expected := sdk.Result{Code: sdk.CodeType(0)}
	res := handler(input.ctx, msg)

	if !ass.Equal(expected.IsOK(), res.IsOK()) {
		t.Errorf("The results should be equal")
	}

}

//unhandled message
func TestNewHandler2(t *testing.T) {

	ass := assert.New(t)

	handler := NewHandler(comAuthkeeper)

	unhandledMsg := id.MsgCreateConnection{}

	err := fmt.Sprintf("Unrecognized commerciodocs message type: %v", unhandledMsg.Type())

	expected := sdk.ErrUnknownRequest(err).Result()

	actual := handler(input.ctx, unhandledMsg)

	if !ass.Equal(expected.Code, actual.Code) {
		t.Errorf("The results should be equal")
	}
}
