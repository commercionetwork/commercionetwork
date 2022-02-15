// DONTCOVER
// nolint
package v3_0_0

import (
	"reflect"
	"testing"

	v220docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var testingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

// var testingSender2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var testingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

const validUUID = "d83422c6-6e79-4a99-9767-fcae46dfa371"
const anotherValidUUID = "49c981c2-a09e-47d2-8814-9373ff64abae"

var testingDocument = types.Document{
	UUID:       validUUID,
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
	EncryptionData: &types.DocumentEncryptionData{
		Keys:          []*types.DocumentEncryptionKey{{Recipient: testingRecipient.String(), Value: "6F7468657276616C7565"}},
		EncryptedData: []string{"content", "content_uri", "metadata.content_uri", "metadata.schema.uri"},
	},
	DoSign: &types.DocumentDoSign{
		StorageURI:     "https://example.com/document/storage",
		SignerInstance: "SignerInstance",
		SdnData: types.SdnData{
			types.SdnDataCommonName,
			types.SdnDataSurname,
			types.SdnDataSurname,
			types.SdnDataGivenName,
			types.SdnDataOrganization,
			types.SdnDataCountry,
		},
		VcrID:              "VcrID",
		CertificateProfile: "CertificateProfile",
	},
	Sender:     testingSender.String(),
	Recipients: []string{testingRecipient.String()},
}

var testingDocumentReceipt = types.DocumentReceipt{
	UUID:         validUUID,
	Sender:       testingSender.String(),
	Recipient:    testingRecipient.String(),
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}

var anotherTestingDocument types.Document
var anotherTestingDocumentReceipt types.DocumentReceipt

var testingv220DocumentRecipients []sdk.AccAddress
var testingv220EncryptionDataKeys []v220docs.DocumentEncryptionKey
var testingv220DocumentReceipt v220docs.DocumentReceipt

var testingv220Document v220docs.Document

var anotherTestingv220Document v220docs.Document
var anotherTestingv220DocumentReceipt v220docs.DocumentReceipt

var invalidTestingv220Document v220docs.Document
var invalidTestingv220DocumentReceipt v220docs.DocumentReceipt

func init() {

	for _, r := range testingDocument.Recipients {
		addr, err := sdk.AccAddressFromBech32(r)
		if err != nil {
			panic("error on addresses for Recipients")
		}
		testingv220DocumentRecipients = append(testingv220DocumentRecipients, addr)
	}

	for _, k := range testingDocument.EncryptionData.Keys {
		address, err := sdk.AccAddressFromBech32(k.Recipient)
		if err != nil {
			panic("error on addresses for EncriptionData Keys")
		}
		testingv220EncryptionDataKeys = append(testingv220EncryptionDataKeys, v220docs.DocumentEncryptionKey{
			Recipient: address,
			Value:     k.Value,
		})
	}

	anotherTestingDocument = testingDocument
	anotherTestingDocument.UUID = anotherValidUUID

	anotherTestingDocumentReceipt = testingDocumentReceipt
	anotherTestingDocumentReceipt.UUID = anotherValidUUID

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
		EncryptionData: &v220docs.DocumentEncryptionData{
			Keys:          testingv220EncryptionDataKeys,
			EncryptedData: testingDocument.EncryptionData.EncryptedData,
		},
		DoSign: &v220docs.DocumentDoSign{
			StorageURI:         testingDocument.DoSign.StorageURI,
			SignerInstance:     testingDocument.DoSign.SignerInstance,
			SdnData:            testingDocument.DoSign.SdnData,
			VcrID:              testingDocument.DoSign.VcrID,
			CertificateProfile: testingDocument.DoSign.CertificateProfile,
		},
	}

	testingv220DocumentReceipt = v220docs.DocumentReceipt{
		UUID:         testingDocumentReceipt.UUID,
		Sender:       testingSender,
		Recipient:    testingRecipient,
		TxHash:       testingDocumentReceipt.TxHash,
		DocumentUUID: testingDocumentReceipt.DocumentUUID,
		Proof:        testingDocumentReceipt.Proof,
	}

	anotherTestingv220Document = testingv220Document
	anotherTestingv220Document.UUID = anotherValidUUID

	anotherTestingv220DocumentReceipt = testingv220DocumentReceipt
	anotherTestingv220DocumentReceipt.UUID = anotherValidUUID

	invalidTestingv220Document = testingv220Document
	invalidTestingv220Document.UUID = "invalid-uuid"

	invalidTestingv220DocumentReceipt = testingv220DocumentReceipt
	invalidTestingv220DocumentReceipt.UUID = "invalid-uuid"
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
			name: "ok",
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

func Test_migrateReceipt(t *testing.T) {
	type args struct {
		receipt v220docs.DocumentReceipt
	}
	tests := []struct {
		name string
		args args
		want *types.DocumentReceipt
	}{
		{
			name: "ok",
			args: args{
				receipt: testingv220DocumentReceipt,
			},
			want: &testingDocumentReceipt,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := migrateReceipt(tt.args.receipt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("migrateReceipt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMigrate(t *testing.T) {
	type args struct {
		oldGenState v220docs.GenesisState
	}
	tests := []struct {
		name string
		args args
		want *types.GenesisState
	}{
		{
			name: "empty",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents:                      []v220docs.Document{},
					Receipts:                       []v220docs.DocumentReceipt{},
					SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
					TrustedMetadataSchemaProposers: []sdk.AccAddress{},
				},
			},
			want: &types.GenesisState{
				Documents: []*types.Document{},
				Receipts:  []*types.DocumentReceipt{},
			},
		},
		{
			name: "one document and corresponding receipt",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						testingv220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						testingv220DocumentReceipt,
					},
					SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
					TrustedMetadataSchemaProposers: []sdk.AccAddress{},
				},
			},
			want: &types.GenesisState{
				Documents: []*types.Document{
					&testingDocument,
				},
				Receipts: []*types.DocumentReceipt{
					&testingDocumentReceipt,
				},
			},
		},
		{
			name: "one document and corresponding receipt, plus one that is invalid so it is not included",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						testingv220Document,
						invalidTestingv220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						testingv220DocumentReceipt,
					},
					SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
					TrustedMetadataSchemaProposers: []sdk.AccAddress{},
				},
			},
			want: &types.GenesisState{
				Documents: []*types.Document{
					&testingDocument,
				},
				Receipts: []*types.DocumentReceipt{
					&testingDocumentReceipt,
				},
			},
		},
		{
			name: "two document and corresponding receipts",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						testingv220Document,
						anotherTestingv220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						testingv220DocumentReceipt,
						anotherTestingv220DocumentReceipt,
					},
					SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
					TrustedMetadataSchemaProposers: []sdk.AccAddress{},
				},
			},
			want: &types.GenesisState{
				Documents: []*types.Document{
					&testingDocument,
					&anotherTestingDocument,
				},
				Receipts: []*types.DocumentReceipt{
					&testingDocumentReceipt,
					&anotherTestingDocumentReceipt,
				},
			},
		},
		{
			name: "one document and corresponding receipt plus a spurious receipt that is not included",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						testingv220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						testingv220DocumentReceipt,
						anotherTestingv220DocumentReceipt,
					},
					SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
					TrustedMetadataSchemaProposers: []sdk.AccAddress{},
				},
			},
			want: &types.GenesisState{
				Documents: []*types.Document{
					&testingDocument,
				},
				Receipts: []*types.DocumentReceipt{
					&testingDocumentReceipt,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := Migrate(tt.args.oldGenState); !reflect.DeepEqual(got, tt.want) {
				require.Equal(t, tt.want, got)
				t.Errorf("Migrate() = %v, want %v", got, tt.want)
			}
		})
	}
}
