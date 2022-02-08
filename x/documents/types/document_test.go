package types

import (
	"strings"
	"testing"
)

func TestDocument_Validate(t *testing.T) {

	tests := []struct {
		name    string
		doc     func() Document
		wantErr bool
	}{
		{
			name: "valid",
			doc: func() Document {
				return ValidDocument
			},
			wantErr: false,
		},
		{
			name: "invalid sender",
			doc: func() Document {
				doc := ValidDocument
				doc.Sender = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "empty recipients",
			doc: func() Document {
				doc := ValidDocument
				doc.Recipients = []string{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid recipients",
			doc: func() Document {
				doc := ValidDocument
				doc.Recipients = []string{"abc"}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid UUID",
			doc: func() Document {
				doc := ValidDocument
				doc.UUID = "abc"
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid metadata",
			doc: func() Document {
				doc := ValidDocument
				doc.Metadata = &DocumentMetadata{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid checksum",
			doc: func() Document {
				doc := ValidDocument
				doc.Checksum = &DocumentChecksum{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid encryption data",
			doc: func() Document {
				doc := ValidDocument
				doc.EncryptionData = &DocumentEncryptionData{}
				return doc
			},
			wantErr: true,
		},
		{
			name: "encryption data does not contain a recipient",
			doc: func() Document {
				doc := ValidDocument
				recipients := append([]string{}, doc.Recipients...)
				doc.Recipients = append(recipients, doc.Sender)
				return doc
			},
			wantErr: true,
		},
		{
			name: "encryption data recipient not included in document recipient list",
			doc: func() Document {
				doc := ValidDocument
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
				doc := ValidDocument
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but empty ContentUri",
			doc: func() Document {
				doc := ValidDocument
				doc.EncryptionData = nil
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but empty MetadataSchemaURI",
			doc: func() Document {
				doc := ValidDocument
				doc.Metadata = &DocumentMetadata{
					ContentURI: ValidDocument.ContentURI,
					Schema:     nil,
				}
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but invalid checksum",
			doc: func() Document {
				doc := ValidDocument
				doc.EncryptionData = nil
				doc.Checksum = nil
				return doc
			},
			wantErr: true,
		},
		{
			name: "invalid SdnData",
			doc: func() Document {
				doc := ValidDocument
				doc.DoSign = &DocumentDoSign{
					StorageURI:         ValidDocument.DoSign.StorageURI,
					SignerInstance:     ValidDocument.DoSign.SignerInstance,
					SdnData:            []string{"planet"},
					VcrID:              ValidDocument.DoSign.VcrID,
					CertificateProfile: ValidDocument.DoSign.CertificateProfile,
				}
				return doc
			},
			wantErr: true,
		},
		{
			name: "content URI not present",
			doc: func() Document {
				doc := ValidDocument
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but empty content uri",
			doc: func() Document {
				doc := ValidDocument
				doc.EncryptionData = nil
				doc.ContentURI = ""
				return doc
			},
			wantErr: true,
		},
		{
			name: "do_sign specified but invalid checksum",
			doc: func() Document {
				doc := ValidDocument
				doc.EncryptionData = nil
				doc.Checksum = nil
				return doc
			},
			wantErr: true,
		},
		{
			name: "violate lenght limits",
			doc: func() Document {
				doc := ValidDocument
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
