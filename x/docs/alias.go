package docs

import (
	"github.com/commercionetwork/commercionetwork/x/docs/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	QuerierRoute = types.QuerierRoute

	MetadataSchemaProposersStoreKey = types.MetadataSchemaProposersStoreKey

	MsgTypeShareDocument = types.MsgTypeShareDocument
)

var (
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgShareDocument = types.NewMsgShareDocument

	SetupTestInput         = keeper.SetupTestInput
	TestingSender          = keeper.TestingSender
	TestingRecipient       = keeper.TestingRecipient
	TestingDocument        = keeper.TestingDocument
	TestingDocumentReceipt = keeper.TestingDocumentReceipt
)

type (
	Keeper                              = keeper.Keeper
	Document                            = types.Document
	DocumentMetadata                    = types.DocumentMetadata
	MetadataSchema                      = types.MetadataSchema
	DocumentChecksum                    = types.DocumentChecksum
	DocumentIds                         = types.DocumentIds
	DocumentReceipt                     = types.DocumentReceipt
	MsgShareDocument                    = types.MsgShareDocument
	MsgSendDocumentReceipt              = types.MsgSendDocumentReceipt
	MsgAddSupportedMetadataSchema       = types.MsgAddSupportedMetadataSchema
	MsgAddTrustedMetadataSchemaProposer = types.MsgAddTrustedMetadataSchemaProposer
)
