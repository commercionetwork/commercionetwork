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
	NewKeeper  = keeper.NewKeeper
	NewHandler = keeper.NewHandler
	NewQuerier = keeper.NewQuerier

	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewMsgShareDocument = types.NewMsgShareDocument

	RegisterInvariants = keeper.RegisterInvariants
)

type (
	Keeper                              = keeper.Keeper
	Document                            = types.Document
	Documents                           = types.Documents
	DocumentMetadata                    = types.DocumentMetadata
	MetadataSchema                      = types.MetadataSchema
	MetadataSchemes                     = types.MetadataSchemes
	DocumentMetadataSchema              = types.DocumentMetadataSchema
	DocumentChecksum                    = types.DocumentChecksum
	DocumentIds                         = types.DocumentIDs
	DocumentReceipt                     = types.DocumentReceipt
	DocumentReceipts                    = types.DocumentReceipts
	MsgShareDocument                    = types.MsgShareDocument
	MsgSendDocumentReceipt              = types.MsgSendDocumentReceipt
	MsgAddSupportedMetadataSchema       = types.MsgAddSupportedMetadataSchema
	MsgAddTrustedMetadataSchemaProposer = types.MsgAddTrustedMetadataSchemaProposer
)
