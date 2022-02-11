package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var recipient1, _ = sdk.AccAddressFromBech32("cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr")
var recipient2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")

var ValidDocument = Document{
	UUID:       "d83422c6-6e79-4a99-9767-fcae46dfa371",
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
		Keys: []*DocumentEncryptionKey{
			{Recipient: recipient1.String(), Value: "6F7468657276616C7565"},
			{Recipient: recipient2.String(), Value: "7F7468657276616C7565"},
		},
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
	Recipients: []string{recipient1.String(), recipient2.String()},
}

var ValidDocumentReceiptRecipient1 = DocumentReceipt{
	UUID:         "8db853ac-5265-4da6-a07a-c52ac8099385",
	Sender:       recipient1.String(),
	Recipient:    sender.String(),
	TxHash:       "txHash",
	DocumentUUID: ValidDocument.UUID,
	Proof:        "proof",
}

var ValidDocumentReceiptRecipient2 = DocumentReceipt{
	UUID:         "bb84a465-6602-43af-9722-7d8a42d81ed8",
	Sender:       recipient2.String(),
	Recipient:    sender.String(),
	TxHash:       "txHash",
	DocumentUUID: ValidDocument.UUID,
	Proof:        "proof",
}

var AnotherValidDocument Document
var AnotherValidDocumentReceipt DocumentReceipt

var InvalidDocument Document
var InvalidDocumentReceipt DocumentReceipt

func init() {
	AnotherValidDocument = ValidDocument
	AnotherValidDocument.UUID = "49c981c2-a09e-47d2-8814-9373ff64abae"

	AnotherValidDocumentReceipt = ValidDocumentReceiptRecipient1
	AnotherValidDocumentReceipt.UUID = "7f4d6197-900a-44af-af22-3a703c568bfe"
	AnotherValidDocumentReceipt.DocumentUUID = AnotherValidDocument.UUID

	InvalidDocument = ValidDocument
	InvalidDocument.UUID = "abc"

	InvalidDocumentReceipt = ValidDocumentReceiptRecipient1
	InvalidDocumentReceipt.UUID = "def"
	InvalidDocumentReceipt.DocumentUUID = InvalidDocument.UUID
}
