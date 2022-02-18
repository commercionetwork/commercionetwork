package v3_0_0

import (
	"reflect"
	"testing"

	v220docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var validDocument = types.ValidDocument
var validDocumentReceipt = types.ValidDocumentReceiptRecipient1
var anotherValidDocument = types.AnotherValidDocument
var anotherDocumentReceipt = types.AnotherValidDocumentReceipt

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
			panic("error on addresses for recipients in document")
		}
		v220DocumentRecipients = append(v220DocumentRecipients, addr)
	}

	for _, k := range validDocument.EncryptionData.Keys {
		address, err := sdk.AccAddressFromBech32(k.Recipient)
		if err != nil {
			panic("error on addresses for encriptionData keys in document")
		}
		v220EncryptionDataKeys = append(v220EncryptionDataKeys, v220docs.DocumentEncryptionKey{
			Recipient: address,
			Value:     k.Value,
		})
	}

	sender, err := sdk.AccAddressFromBech32(validDocument.Sender)
	if err != nil {
		panic("error on address for sender in receipt")
	}
	recipient, err := sdk.AccAddressFromBech32(validDocumentReceipt.Sender)
	if err != nil {
		panic("error on address for recipient in receipt")
	}

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
		Sender:       recipient,
		Recipient:    sender,
		TxHash:       validDocumentReceipt.TxHash,
		DocumentUUID: validDocumentReceipt.DocumentUUID,
		Proof:        validDocumentReceipt.Proof,
	}

	anotherV220Document = v220docs.Document{
		Sender:     sender,
		Recipients: v220DocumentRecipients,
		UUID:       anotherValidDocument.UUID,
		Metadata: v220docs.DocumentMetadata{
			ContentURI: validDocument.Metadata.ContentURI,
			Schema: &v220docs.DocumentMetadataSchema{
				URI:     validDocument.Metadata.Schema.URI,
				Version: validDocument.Metadata.Schema.Version,
			},
		},
		ContentURI: validDocument.ContentURI,
		Checksum: &v220docs.DocumentChecksum{
			Value:     "825b5f1e6f9fe03eac07d27b164af1a2",
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
						v220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						v220DocumentReceipt,
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
		{
			name: "one document and corresponding receipt, plus one that is invalid so it is not included",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						v220Document,
						invalidV220Document,
					},
					Receipts: []v220docs.DocumentReceipt{
						v220DocumentReceipt,
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
		{
			name: "two document and corresponding receipts",
			args: args{
				oldGenState: v220docs.GenesisState{
					Documents: []v220docs.Document{
						v220Document,
						anotherV220Document,
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
					&anotherValidDocument,
				},
				Receipts: []*types.DocumentReceipt{
					&validDocumentReceipt,
					&anotherDocumentReceipt,
				},
			},
		},
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
