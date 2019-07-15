package commercioauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/assert"
	"testing"
)

var cAuthkeeper = Keeper{
	accountKeeper: input.accKeeper,
	cdc:           input.cdc,
}

//handled query
func TestNewQuerier(t *testing.T) {

	ass := assert.New(t)
	querier := NewQuerier(cAuthkeeper)

	path := []string{"account", "list"}

	var request abci.RequestQuery

	var expectedBytes = []byte{0x6e, 0x75, 0x6c, 0x6c}

	actual, _ := querier(input.ctx, path, request)

	if !ass.Equal(expectedBytes, actual) {
		t.Errorf("The two byte arrays should be equals")
	}

}

//unhandled request
func TestNewQuerier2(t *testing.T) {

	ass := assert.New(t)
	querier := NewQuerier(cAuthkeeper)

	path := []string{"acc", "list"}

	var request abci.RequestQuery

	var expectedError = sdk.ErrUnknownRequest("Unknown commerciodocs query endpoint")

	_, actualErr := querier(input.ctx, path, request)

	if !ass.Equal(expectedError.Result(), actualErr.Result()) {
		t.Errorf("The two errors should be equals")
	}
}
