package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var testUtils = setupTestInput()

//TEST VARS
var testAddress = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var testOwner, _ = sdk.AccAddressFromBech32(testAddress)
var testOwnerIdentity = types.Did("newReader")
var testIdentityRef = "ddo-reference"
var testReference = "testReference"
var testMetadata = "testMetadata"
var testRecipient = types.Did("recipient")
