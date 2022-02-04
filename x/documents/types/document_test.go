package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var recipient, _ = sdk.AccAddressFromBech32("cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr")

const validUUID = "d83422c6-6e79-4a99-9767-fcae46dfa371"
const anotherValidUUID = "49c981c2-a09e-47d2-8814-9373ff64abae"

var validDocument = Document{
	UUID:       validUUID,
	ContentURI: "https://example.com/document",
	Metadata: &DocumentMetadata{
		ContentURI: "https://example.com/document/metadata",
		Schema: &DocumentMetadataSchema{
			URI:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
	},
	Checksum: &DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
	EncryptionData: &DocumentEncryptionData{
		Keys:          []*DocumentEncryptionKey{{Recipient: recipient.String(), Value: "6F7468657276616C7565"}},
		EncryptedData: []string{"content", "content_uri", "metadata.content_uri", "metadata.schema.uri"},
	},
	DoSign: &DocumentDoSign{
		StorageURI:     "https://example.com/document/storage",
		SignerInstance: "SignerInstance",
		SdnData: SdnData{
			SdnDataCommonName,
			SdnDataSurname,
			SdnDataSurname,
			SdnDataGivenName,
			SdnDataOrganization,
			SdnDataCountry,
		},
		VcrID:              "VcrID",
		CertificateProfile: "CertificateProfile",
	},
	Sender:     sender.String(),
	Recipients: []string{recipient.String()},
}

var validDocumentReceipt = DocumentReceipt{
	UUID:         validUUID,
	Sender:       sender.String(),
	Recipient:    recipient.String(),
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}

func TestDocument_Validate(t *testing.T) {

	tests := []struct {
		name    string
		doc     func() Document
		wantErr bool
	}{
		{
			name: "valid",
			doc: func() Document {
				return validDocument
			},
			wantErr: false,
		},
		{
			name: "empty sender",
			doc: func() Document {
				doc := validDocument
				doc.Sender = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "empty recipients",
			doc: func() Document {
				doc := validDocument
				doc.Recipients = []string{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "empty recipients",
			doc: func() Document {
				doc := validDocument
				doc.Recipients = []string{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid recipients",
			doc: func() Document {
				doc := validDocument
				doc.Recipients = []string{"abc"}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid UUID",
			doc: func() Document {
				doc := validDocument
				doc.UUID = "abc"
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid metadata",
			doc: func() Document {
				doc := validDocument
				doc.Metadata = &DocumentMetadata{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid checksum",
			doc: func() Document {
				doc := validDocument
				doc.Checksum = &DocumentChecksum{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid encryption data",
			doc: func() Document {
				doc := validDocument
				doc.EncryptionData = &DocumentEncryptionData{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "encryption data does not contain a recipient",
			doc: func() Document {
				doc := validDocument
				doc.Recipients = append(doc.Recipients, doc.Sender)
				return doc
			},
			wantErr: true,
		},
		// {
		// 	name: "check that there are no spurious encryption data recipients not present in the document recipient list",
		// 	wantErr: true,
		// },
		// {
		// 	name: "Check that the `encrypted_data' field name is actually present in doc",
		// 	wantErr: true,
		// },
		// {
		// 	name: "do_sign",
		// 	wantErr: true,
		// },
		{
			name: "violate lenght limits",
			doc: func() Document {
				doc := validDocument
				doc.ContentURI = strings.Repeat("a", 513)
				return doc
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := tt.doc()
			if err := doc.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Document.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
