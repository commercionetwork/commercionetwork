package types

import (
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TEST VARS
var addr = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var sender, _ = sdk.AccAddressFromBech32(addr)
var receiver, _ = sdk.AccAddressFromBech32(addr)
var validChecksum = types.DocumentChecksum{
	Value:     "testValue",
	Algorithm: SHA256,
}
var validMetadataSchema = types.DocumentMetadataSchema{
	Uri:     "http://www.contentUri.com",
	Version: "test",
}
var validDocumentMetadata = types.DocumentMetadata{
	ContentUri: "http://www.contentUri.com",
	Schema:     validMetadataSchema,
	Proof:      "proof",
}

var invalidChecksum = types.DocumentChecksum{
	Value:     "",
	Algorithm: "",
}
var invalidMetadataSchema = types.DocumentMetadataSchema{
	Uri:     "",
	Version: "",
}
var invalidDocumentMetadata = types.DocumentMetadata{
	ContentUri: "",
	Schema:     invalidMetadataSchema,
	Proof:      "",
}

var validMsg = MsgShareDocument{types.Document{
	Sender:     sender,
	Recipient:  receiver,
	Uuid:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentUri: "http://www.contentUri.com",
	Metadata:   validDocumentMetadata,
	Checksum:   validChecksum,
}}

var invalidMsg = MsgShareDocument{types.Document{
	Sender:     sender,
	Recipient:  receiver,
	Uuid:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentUri: "http://www.contentUri.com",
	Metadata: types.DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: types.DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	},
	Checksum: types.DocumentChecksum{
		Value:     "testValue",
		Algorithm: SHA256,
	},
}}

func TestMsgShareDocument_Route(t *testing.T) {
	actual := validMsg.Route()

}
