// DONTCOVER
// nolint
package v3_0_0

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	v220docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

var testingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

// var testingSender2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var testingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var testingDocument = types.Document{
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
	Sender:     testingSender.String(),
	Recipients: []string{testingRecipient.String()},
}

// var testingDocumentReceipt = types.DocumentReceipt{
// 	UUID:         "testing-document-receipt-uuid",
// 	Sender:       testingSender.String(),
// 	Recipient:    testingRecipient.String(),
// 	TxHash:       "txHash",
// 	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
// 	Proof:        "proof",
// }

var testingv220Document v220docs.Document
var testingv220DocumentRecipients []sdk.AccAddress

func init() {

	for _, r := range testingDocument.Recipients {
		addr, err := sdk.AccAddressFromBech32(r)
		if err != nil {
			panic("error on addresses")
		}
		testingv220DocumentRecipients = append(testingv220DocumentRecipients, addr)
	}

	testingv220Document = v220docs.Document{
		Sender:     testingSender,
		Recipients: testingv220DocumentRecipients,
		UUID:       testingDocument.UUID,
		Metadata: v220docs.DocumentMetadata{
			ContentURI: testingDocument.Metadata.ContentURI,
			// SchemaType: "", // not converted
			Schema: &v220docs.DocumentMetadataSchema{
				URI:     testingDocument.Metadata.Schema.URI,
				Version: testingDocument.Metadata.Schema.Version,
			},
		},
		ContentURI: testingDocument.ContentURI,
		Checksum: &v220docs.DocumentChecksum{
			Value:     testingDocument.Checksum.Value,
			Algorithm: testingDocument.Checksum.Algorithm,
		},
		// EncryptionData: &v220docs.DocumentEncryptionData{},
		// DoSign: &v220docs.DocumentDoSign{},
	}
}

func Test_migrateDocument(t *testing.T) {
	type args struct {
		doc v220docs.Document
	}
	tests := []struct {
		name string
		args args
		want *types.Document
	}{
		{
			name: "",
			args: args{
				doc: testingv220Document,
			},
			want: &testingDocument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := migrateDocument(tt.args.doc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("migrateDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}
