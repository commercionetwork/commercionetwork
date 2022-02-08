package v3_0_0

import (
	"reflect"
	"testing"

	v220docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var recipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var validDocument = types.Document{
	UUID:       "d83422c6-6e79-4a99-9767-fcae46dfa371",
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
		Keys:          []*types.DocumentEncryptionKey{{Recipient: recipient.String(), Value: "6F7468657276616C7565"}},
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
	Sender:     sender.String(),
	Recipients: []string{recipient.String()},
}

var validDocumentReceipt = types.DocumentReceipt{
	UUID:         "9680b9f8-d3a2-49c4-b32e-517632e6e67d",
	Sender:       sender.String(),
	Recipient:    recipient.String(),
	TxHash:       "txHash",
	DocumentUUID: validDocument.UUID,
	Proof:        "proof",
}

var anotherValidDocument types.Document
var anotherDocumentReceipt types.DocumentReceipt

var v220DocumentRecipients []sdk.AccAddress
var v220EncryptionDataKeys []v220docs.DocumentEncryptionKey
var v220DocumentReceipt v220docs.DocumentReceipt

var v220Document v220docs.Document

var anotherV220Document v220docs.Document
var anotherV220DocumentReceipt v220docs.DocumentReceipt

var invalidV220Document v220docs.Document
var invalidV220DocumentReceipt v220docs.DocumentReceipt

func init() {

	for _, r := range validDocument.Recipients {
		addr, err := sdk.AccAddressFromBech32(r)
		if err != nil {
			panic("error on addresses for Recipients")
		}
		v220DocumentRecipients = append(v220DocumentRecipients, addr)
	}

	for _, k := range validDocument.EncryptionData.Keys {
		address, err := sdk.AccAddressFromBech32(k.Recipient)
		if err != nil {
			panic("error on addresses for EncriptionData Keys")
		}
		v220EncryptionDataKeys = append(v220EncryptionDataKeys, v220docs.DocumentEncryptionKey{
			Recipient: address,
			Value:     k.Value,
		})
	}

	anotherValidDocument = validDocument
	anotherValidDocument.UUID = "49c981c2-a09e-47d2-8814-9373ff64abae"

	anotherDocumentReceipt = validDocumentReceipt
	anotherDocumentReceipt.UUID = "7f4d6197-900a-44af-af22-3a703c568bfe"

	v220Document = v220docs.Document{
		Sender:     sender,
		Recipients: v220DocumentRecipients,
		UUID:       validDocument.UUID,
		Metadata: v220docs.DocumentMetadata{
			ContentURI: validDocument.Metadata.ContentURI,
			// SchemaType: "", // not converted
			Schema: &v220docs.DocumentMetadataSchema{
				URI:     validDocument.Metadata.Schema.URI,
				Version: validDocument.Metadata.Schema.Version,
			},
		},
		ContentURI: validDocument.ContentURI,
		Checksum: &v220docs.DocumentChecksum{
			Value:     validDocument.Checksum.Value,
			Algorithm: validDocument.Checksum.Algorithm,
		},
		EncryptionData: &v220docs.DocumentEncryptionData{
			Keys:          v220EncryptionDataKeys,
			EncryptedData: validDocument.EncryptionData.EncryptedData,
		},
		DoSign: &v220docs.DocumentDoSign{
			StorageURI:         validDocument.DoSign.StorageURI,
			SignerInstance:     validDocument.DoSign.SignerInstance,
			SdnData:            validDocument.DoSign.SdnData,
			VcrID:              validDocument.DoSign.VcrID,
			CertificateProfile: validDocument.DoSign.CertificateProfile,
		},
	}

	v220DocumentReceipt = v220docs.DocumentReceipt{
		UUID:         validDocumentReceipt.UUID,
		Sender:       sender,
		Recipient:    recipient,
		TxHash:       validDocumentReceipt.TxHash,
		DocumentUUID: validDocumentReceipt.DocumentUUID,
		Proof:        validDocumentReceipt.Proof,
	}

	anotherV220Document = v220Document
	anotherV220Document.UUID = anotherValidDocument.UUID

	anotherV220DocumentReceipt = v220DocumentReceipt
	anotherV220DocumentReceipt.UUID = anotherDocumentReceipt.UUID
	anotherV220DocumentReceipt.DocumentUUID = anotherV220Document.UUID

	invalidV220Document = v220Document
	invalidV220Document.UUID = "invalid-uuid"

	invalidV220DocumentReceipt = v220DocumentReceipt
	invalidV220DocumentReceipt.UUID = "invalid-uuid"
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
				doc: v220Document,
			},
			want: &validDocument,
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
				receipt: v220DocumentReceipt,
			},
			want: &validDocumentReceipt,
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
		// {
		// 	name: "empty",
		// 	args: args{
		// 		oldGenState: v220docs.GenesisState{
		// 			Documents:                      []v220docs.Document{},
		// 			Receipts:                       []v220docs.DocumentReceipt{},
		// 			SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
		// 			TrustedMetadataSchemaProposers: []sdk.AccAddress{},
		// 		},
		// 	},
		// 	want: &types.GenesisState{
		// 		Documents: []*types.Document{},
		// 		Receipts:  []*types.DocumentReceipt{},
		// 	},
		// },
		// {
		// 	name: "one document and corresponding receipt",
		// 	args: args{
		// 		oldGenState: v220docs.GenesisState{
		// 			Documents: []v220docs.Document{
		// 				v220Document,
		// 			},
		// 			Receipts: []v220docs.DocumentReceipt{
		// 				v220DocumentReceipt,
		// 			},
		// 			SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
		// 			TrustedMetadataSchemaProposers: []sdk.AccAddress{},
		// 		},
		// 	},
		// 	want: &types.GenesisState{
		// 		Documents: []*types.Document{
		// 			&validDocument,
		// 		},
		// 		Receipts: []*types.DocumentReceipt{
		// 			&validDocumentReceipt,
		// 		},
		// 	},
		// },
		// {
		// 	name: "one document and corresponding receipt, plus one that is invalid so it is not included",
		// 	args: args{
		// 		oldGenState: v220docs.GenesisState{
		// 			Documents: []v220docs.Document{
		// 				v220Document,
		// 				invalidV220Document,
		// 			},
		// 			Receipts: []v220docs.DocumentReceipt{
		// 				v220DocumentReceipt,
		// 			},
		// 			SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
		// 			TrustedMetadataSchemaProposers: []sdk.AccAddress{},
		// 		},
		// 	},
		// 	want: &types.GenesisState{
		// 		Documents: []*types.Document{
		// 			&validDocument,
		// 		},
		// 		Receipts: []*types.DocumentReceipt{
		// 			&validDocumentReceipt,
		// 		},
		// 	},
		// },
		// {
		// 	name: "two document and corresponding receipts",
		// 	args: args{
		// 		oldGenState: v220docs.GenesisState{
		// 			Documents: []v220docs.Document{
		// 				v220Document,
		// 				anotherV220Document,
		// 			},
		// 			Receipts: []v220docs.DocumentReceipt{
		// 				v220DocumentReceipt,
		// 				anotherV220DocumentReceipt,
		// 			},
		// 			SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
		// 			TrustedMetadataSchemaProposers: []sdk.AccAddress{},
		// 		},
		// 	},
		// 	want: &types.GenesisState{
		// 		Documents: []*types.Document{
		// 			&validDocument,
		// 			&anotherValidDocument,
		// 		},
		// 		Receipts: []*types.DocumentReceipt{
		// 			&validDocumentReceipt,
		// 			&anotherDocumentReceipt,
		// 		},
		// 	},
		// },
		{
			name: "one document and corresponding receipt plus a spurious receipt that is not included",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						v220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						v220DocumentReceipt,
						anotherV220DocumentReceipt,
					},
					SupportedMetadataSchemes:       []v220docs.MetadataSchema{},
					TrustedMetadataSchemaProposers: []sdk.AccAddress{},
				},
			},
			want: &types.GenesisState{
				Documents: []*types.Document{
					&validDocument,
				},
				Receipts: []*types.DocumentReceipt{
					&validDocumentReceipt,
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
