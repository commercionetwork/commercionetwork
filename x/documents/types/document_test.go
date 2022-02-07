package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var recipient, _ = sdk.AccAddressFromBech32("cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr")

const validDocumentUUID = "d83422c6-6e79-4a99-9767-fcae46dfa371"
const anotherValidDocumentUUID = "49c981c2-a09e-47d2-8814-9373ff64abae"

const validReceiptUUID = "8db853ac-5265-4da6-a07a-c52ac8099385"

var validDocument = Document{
	UUID:       validDocumentUUID,
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
	UUID:         validReceiptUUID,
	Sender:       sender.String(),
	Recipient:    recipient.String(),
	TxHash:       "txHash",
	DocumentUUID: validDocumentUUID,
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
			name: "invalid sender",
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
				recipients := append([]string{}, doc.Recipients...)
				doc.Recipients = append(recipients, doc.Sender)
				return doc
			},
			wantErr: true,
		},
		{
			name: "encryption data recipient not included in document recipient list",
			doc: func() Document {
				doc := validDocument
				encryptionDataKeys := append([]*DocumentEncryptionKey{}, doc.EncryptionData.Keys...)
				encryptionDataKeys = append(encryptionDataKeys, &DocumentEncryptionKey{
					Recipient: doc.Sender,
					Value:     "6F7468657276616C7565",
				})

				encryptionData := DocumentEncryptionData{
					Keys:          encryptionDataKeys,
					EncryptedData: validDocumentEncryptionData.EncryptedData,
				}

				doc.EncryptionData = &encryptionData
				return doc
			},
			wantErr: true,
		},
		{
			name: "content URI not present",
			doc: func() Document {
				doc := validDocument
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but empty ContentUri",
			doc: func() Document {
				doc := validDocument
				doc.EncryptionData = nil
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but empty MetadataSchemaURI",
			doc: func() Document {
				doc := validDocument
				doc.Metadata = &DocumentMetadata{
					ContentURI: validDocument.ContentURI,
					Schema:     nil,
				}
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but invalid checksum",
			doc: func() Document {
				doc := validDocument
				doc.EncryptionData = nil
				doc.Checksum = nil
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid SdnData",
			doc: func() Document {
				doc := validDocument
				doc.DoSign = &DocumentDoSign{
					StorageURI:         validDocument.DoSign.StorageURI,
					SignerInstance:     validDocument.DoSign.SignerInstance,
					SdnData:            []string{"planet"},
					VcrID:              validDocument.DoSign.VcrID,
					CertificateProfile: validDocument.DoSign.CertificateProfile,
				}
				return doc
			},
			wantErr: true,
		},
		{
			name: "content URI not present",
			doc: func() Document {
				doc := validDocument
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but empty content uri",
			doc: func() Document {
				doc := validDocument
				doc.EncryptionData = nil
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but invalid checksum",
			doc: func() Document {
				doc := validDocument
				doc.EncryptionData = nil
				doc.Checksum = nil
				return doc
			},
			wantErr: true,
		},
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
