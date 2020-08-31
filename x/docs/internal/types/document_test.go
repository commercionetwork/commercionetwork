package types

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"

	"github.com/commercionetwork/commercionetwork/x/common/types"
)

func TestDocument_Equals_NilValues(t *testing.T) {
	document := Document{
		UUID: "uuid",
		Metadata: DocumentMetadata{
			ContentURI: "document_metadata_content_uri",
			SchemaType: "document_metadata_schema_type",
		},
		ContentURI:     "",
		Checksum:       nil,
		EncryptionData: nil,
	}
	require.True(t, document.Equals(document))
}

func Test_validateUUID(t *testing.T) {
	tests := []struct {
		name    string
		UUID    string
		badUUID bool
	}{
		{
			"empty string",
			"",
			true,
		},
		{
			"a well-formed UUID",
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			false,
		},
		{
			"a seemingly well-formed UUID, with the last character removed",
			"6ba7b810-9dad-11d1-80b4-00c04fd430c",
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val := validateUUID(tt.UUID)
			if tt.badUUID {
				require.False(t, val, "got true")
			} else {
				require.True(t, val, "got false")
			}
		})
	}
}

func TestDocument_Equals(t *testing.T) {
	sender, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	recipient, _ := sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")
	tests := []struct {
		name      string
		other     Document
		us        Document
		different bool
	}{
		{
			"two empty documents",
			Document{},
			Document{},
			false,
		},
		{
			"two equal documents",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				UUID:           "ac33043b-5cb4-4645-a3f9-819140847252",
				Checksum:       &DocumentChecksum{},
				EncryptionData: &DocumentEncryptionData{},
			},
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				UUID:           "ac33043b-5cb4-4645-a3f9-819140847252",
				Checksum:       &DocumentChecksum{},
				EncryptionData: &DocumentEncryptionData{},
			},
			false,
		},
		{
			"two identical documents, except the UUID",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				UUID:           "ac33043b-5cb4-4645-a3f9-81914084725",
				Checksum:       &DocumentChecksum{},
				EncryptionData: &DocumentEncryptionData{},
			},
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				UUID:           "ac33043b-5cb4-4645-a3f9-819140847252",
				Checksum:       &DocumentChecksum{},
				EncryptionData: &DocumentEncryptionData{},
			},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if !tt.different {
				require.True(t, tt.us.Equals(tt.other))
			} else {
				require.False(t, tt.us.Equals(tt.other))
			}
		})
	}
}

func TestDocument_Validate(t *testing.T) {
	sender, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	recipient, _ := sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")
	anotherRecipient, _ := sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

	tests := []struct {
		name        string
		doc         Document
		expectedErr error
	}{
		{
			"a good document",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
			},
			nil,
		},
		{
			"no sender",
			Document{
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, ""),
		},
		{
			"no recipients",
			Document{
				Sender: sender,
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Recipients cannot be empty"),
		},
		{
			"no uuid",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid document UUID: "),
		},
		{
			"a good document with some encrypted data inside",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				EncryptionData: &DocumentEncryptionData{
					Keys: []DocumentEncryptionKey{
						{
							Recipient: recipient,
							Value:     "6b6579",
						},
					},
					EncryptedData: []string{"content"},
				},
			},
			nil,
		},
		{
			"a good document whom encrypted data recipient isn't contained in the recipients",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				EncryptionData: &DocumentEncryptionData{
					Keys: []DocumentEncryptionKey{
						{
							Recipient: recipient,
							Value:     "6b6579",
						},
						{
							Recipient: anotherRecipient,
							Value:     "6b6579",
						},
					},
					EncryptedData: []string{"content"},
				},
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf(
				"%s is a recipient inside encryption data but not inside the message",
				anotherRecipient.String(),
			)),
		},
		{
			"a good document whom encrypted data recipient isn't in the document recipient",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					anotherRecipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				EncryptionData: &DocumentEncryptionData{
					Keys: []DocumentEncryptionKey{
						{
							Recipient: recipient,
							Value:     "6b6579",
						},
					},
					EncryptedData: []string{"content"},
				},
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("%s is a recipient inside the document but not in the encryption data", anotherRecipient.String())),
		},
		{
			"a good document whom encrypted data is content_uri, and the corresponding field isn't available",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				EncryptionData: &DocumentEncryptionData{
					Keys: []DocumentEncryptionKey{
						{
							Recipient: recipient,
							Value:     "6b6579",
						},
					},
					EncryptedData: []string{"content_uri"},
				},
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("field \"%s\" not present in document, but marked as encrypted", "content_uri"),
			),
		},
		{
			"a good document whom encrypted data is metadata.schema.uri, and the corresponding field isn't available",
			Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				EncryptionData: &DocumentEncryptionData{
					Keys: []DocumentEncryptionKey{
						{
							Recipient: recipient,
							Value:     "6b6579",
						},
					},
					EncryptedData: []string{"metadata.schema.uri"},
				},
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("field \"%s\" not present in document, but marked as encrypted", "metadata.schema.uri"),
			),
		},
		{
			name: "good document that has do_sign but no checksum",
			doc: Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				DoSign: &DocumentDoSign{
					StorageURI: "theuri",
				},
			},
			expectedErr: sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("field \"%s\" not present in document, but required when using do_sign", "checksum"),
			),
		},
		{
			name: "good document that has do_sign but no content_uri",
			doc: Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Checksum: &DocumentChecksum{
					Value:     "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8",
					Algorithm: "sha-1",
				},
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				DoSign: &DocumentDoSign{
					StorageURI: "theuri",
				},
			},
			expectedErr: sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("field \"%s\" not present in document, but required when using do_sign", "content_uri"),
			),
		}, {
			name: "good document that has do_sign but invalid do_sign sdndata",
			doc: Document{
				Sender: sender,
				Recipients: types.Addresses{
					recipient,
				},
				Checksum: &DocumentChecksum{
					Value:     "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8",
					Algorithm: "sha-1",
				},
				ContentURI: "theContentUri",
				Metadata: DocumentMetadata{
					ContentURI: "content_uri",
					SchemaType: "a schema type",
				},
				UUID: "ac33043b-5cb4-4645-a3f9-819140847252",
				DoSign: &DocumentDoSign{
					StorageURI: "theuri",
					SdnData: SdnData{
						"invalid",
					},
				},
			},
			expectedErr: sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				"sdn_data value \"invalid\" is not supported",
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.doc.Validate()
			if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreateDoc_DoSign(t *testing.T) {
	baseDocument := Document{
		UUID: "uuid",
		Metadata: DocumentMetadata{
			ContentURI: "document_metadata_content_uri",
			SchemaType: "document_metadata_schema_type",
		},
		ContentURI:     "",
		Checksum:       nil,
		EncryptionData: nil,
	}

	tests := []struct {
		name           string
		documentDoSign *DocumentDoSign
		expectedResult string
	}{
		{
			"no do sign (null)",
			nil,
			"{\"sender\":\"\",\"recipients\":null,\"uuid\":\"uuid\",\"metadata\":{\"content_uri\":\"document_metadata_content_uri\",\"schema_type\":\"document_metadata_schema_type\"}}",
		},
		{
			"some data but empty sdn data",
			&DocumentDoSign{
				StorageURI:         "abc",
				SignerInstance:     "abc",
				SdnData:            SdnData{},
				VcrID:              "abc",
				CertificateProfile: "abc",
			},
			"{\"sender\":\"\",\"recipients\":null,\"uuid\":\"uuid\",\"metadata\":{\"content_uri\":\"document_metadata_content_uri\",\"schema_type\":\"document_metadata_schema_type\"},\"do_sign\":{\"storage_uri\":\"abc\",\"signer_instance\":\"abc\",\"sdn_data\":[],\"vcr_id\":\"abc\",\"certificate_profile\":\"abc\"}}",
		},
		{
			"all data",
			&DocumentDoSign{
				StorageURI:     "abc",
				SignerInstance: "abc",
				SdnData: SdnData{
					"common_name",
					"surname",
					"serial_number",
					"given_name",
					"organization",
					"country",
				},
				VcrID:              "abc",
				CertificateProfile: "abc",
			},
			"{\"sender\":\"\",\"recipients\":null,\"uuid\":\"uuid\",\"metadata\":{\"content_uri\":\"document_metadata_content_uri\",\"schema_type\":\"document_metadata_schema_type\"},\"do_sign\":{\"storage_uri\":\"abc\",\"signer_instance\":\"abc\",\"sdn_data\":[\"common_name\",\"surname\",\"serial_number\",\"given_name\",\"organization\",\"country\"],\"vcr_id\":\"abc\",\"certificate_profile\":\"abc\"}}",
		},
	}

	cdc := amino.NewCodec()

	for _, tt := range tests {
		tt := tt
		baseDoc := baseDocument

		t.Run(tt.name, func(t *testing.T) {
			baseDoc.DoSign = tt.documentDoSign
			json, err := cdc.MarshalJSON(baseDoc)
			require.NoError(t, err)
			require.Equal(t, tt.expectedResult, string(json))
		})
	}
}
