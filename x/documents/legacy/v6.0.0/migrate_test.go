package v600_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	v300 "github.com/commercionetwork/commercionetwork/x/documents/legacy/v3.0.0"
	v600 "github.com/commercionetwork/commercionetwork/x/documents/legacy/v6.0.0"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/stretchr/testify/require"
)

var (
	sender, _     = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	recipient1, _ = sdk.AccAddressFromBech32("cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr")
	recipient2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")

	SdnDataCommonName   = "common_name"
	SdnDataSurname      = "surname"
	SdnDataSerialNumber = "serial_number"
	SdnDataGivenName    = "given_name"
	SdnDataOrganization = "organization"
	SdnDataCountry      = "country"

	oldDocument = v300.Document{
		UUID:       "d83422c6-6e79-4a99-9767-fcae46dfa371",
		ContentURI: "https://example.com/document",
		Metadata: &v300.DocumentMetadata{
			ContentURI: "https://example.com/document/metadata",
			Schema: &v300.DocumentMetadataSchema{
				URI:     "https://example.com/document/metadata/schema",
				Version: "1.0.0",
			},
		},
		Checksum: &v300.DocumentChecksum{
			Value:     "93dfcaf3d923ec47edb8580667473987",
			Algorithm: "md5",
		},
		EncryptionData: &v300.DocumentEncryptionData{
			Keys: []*v300.DocumentEncryptionKey{
				{Recipient: recipient1.String(), Value: "6F7468657276616C7565"},
				{Recipient: recipient2.String(), Value: "7F7468657276616C7565"},
			},
			EncryptedData: []string{"content", "content_uri", "metadata.content_uri", "metadata.schema.uri"},
		},
		DoSign: &v300.DocumentDoSign{
			StorageURI:     "https://example.com/document/storage",
			SignerInstance: "SignerInstance",
			SdnData: types.SdnData{
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

	anotherOldDocument = v300.Document{
		UUID:       "49c981c2-a09e-47d2-8814-9373ff64abae",
		ContentURI: "https://example.com/document",
		Metadata: &v300.DocumentMetadata{
			ContentURI: "https://example.com/document/metadata",
			Schema: &v300.DocumentMetadataSchema{
				URI:     "https://example.com/document/metadata/schema",
				Version: "1.0.0",
			},
		},
		Sender:     sender.String(),
		Recipients: []string{recipient1.String(), recipient2.String()},
	}

	oldDocumentReceipt = v300.DocumentReceipt{
		UUID:         "8db853ac-5265-4da6-a07a-c52ac8099385",
		Sender:       recipient1.String(),
		Recipient:    sender.String(),
		TxHash:       "txHash",
		DocumentUUID: oldDocument.UUID,
		Proof:        "proof",
	}
)

func TestMigrateStore(t *testing.T) {
	db := dbm.NewMemDB()
	storeKey := sdk.NewKVStoreKey(v300.StoreKey)

	// Create a multi-store and mount the store key
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, ms.LoadLatestVersion())

	ctx := sdk.NewContext(ms, tmproto.Header{}, false, nil)
	cdc := codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
	store := ctx.KVStore(storeKey)

	store.Set([]byte(v300.DocumentStorePrefix+oldDocument.UUID), cdc.MustMarshal(&oldDocument))
	store.Set([]byte(v300.DocumentStorePrefix+anotherOldDocument.UUID), cdc.MustMarshal(&anotherOldDocument))
	store.Set([]byte(types.ReceiptsStorePrefix+oldDocumentReceipt.UUID), cdc.MustMarshal(&oldDocumentReceipt))

	err := v600.MigrateStore(ctx, storeKey, cdc)
	require.NoError(t, err)

	// Check if the document was migrated correctly
	var newDocument types.Document
	bz := store.Get([]byte(v300.DocumentStorePrefix + oldDocument.UUID))
	require.NotNil(t, bz)
	cdc.MustUnmarshal(bz, &newDocument)

	expectedTimestamp := time.Unix(0, 0).UTC()
	expectedDocument := types.Document{
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
			Keys: []*types.DocumentEncryptionKey{
				{Recipient: recipient1.String(), Value: "6F7468657276616C7565"},
				{Recipient: recipient2.String(), Value: "7F7468657276616C7565"},
			},
			EncryptedData: []string{"content", "content_uri", "metadata.content_uri", "metadata.schema.uri"},
		},
		DoSign: &types.DocumentDoSign{
			StorageURI:     "https://example.com/document/storage",
			SignerInstance: "SignerInstance",
			SdnData: types.SdnData{
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
		Timestamp:  &expectedTimestamp,
	}

	require.Equal(t, expectedDocument, newDocument)

	// Check if the document receipt was migrated correctly
	var newDocumentReceipt types.DocumentReceipt
	bz = store.Get([]byte(types.ReceiptsStorePrefix + oldDocumentReceipt.UUID))
	require.NotNil(t, bz)
	cdc.MustUnmarshal(bz, &newDocumentReceipt)

	expectedReceipt := types.DocumentReceipt{
		UUID:         "8db853ac-5265-4da6-a07a-c52ac8099385",
		Sender:       recipient1.String(),
		Recipient:    sender.String(),
		TxHash:       "txHash",
		DocumentUUID: oldDocument.UUID,
		Proof:        "proof",
		Timestamp:    &expectedTimestamp,
	}

	require.Equal(t, expectedReceipt, newDocumentReceipt)
}
