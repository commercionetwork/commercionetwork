package commerciodocs

import (
	"commercio-network/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var address = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var owner, _ = sdk.AccAddressFromBech32(address)
var ownerIdentity = types.Did("newReader")
var reference = "reference"
var metadata = "metadata"
var recipient = types.Did("recipient")

var msgStore = MsgStoreDocument{
	Owner:     owner,
	Identity:  ownerIdentity,
	Reference: reference,
	Metadata:  metadata,
}

var msgShare = MsgShareDocument{
	Owner:     owner,
	Sender:    ownerIdentity,
	Receiver:  recipient,
	Reference: reference,
}
