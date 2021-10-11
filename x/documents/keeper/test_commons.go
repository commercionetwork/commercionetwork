package keeper

import (
	/*"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"*/
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	/*"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"*/)

// Testing variables

var TestingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestingSender2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var TestingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var TestingDocument = types.Document{
	UUID:       "test-document-uuid",
	ContentURI: "https://example.com/document",
	Metadata: &types.DocumentMetadata{
		ContentURI: "https://example.com/document/metadata",
		Schema: &types.DocumentMetadataSchema{
			URI:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
	},
	Checksum: &types.DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
	Sender:     TestingSender.String(),
	Recipients: append([]string{}, TestingRecipient.String()),
}

var TestingDocumentReceipt = types.DocumentReceipt{
	UUID:         "testing-document-receipt-uuid",
	Sender:       TestingSender.String(),
	Recipient:    TestingRecipient.String(),
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}
