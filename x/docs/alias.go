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
	DocumentMetadata                    = types.DocumentMetadata
	MetadataSchema                      = types.MetadataSchema
	DocumentMetadataSchema              = types.DocumentMetadataSchema
	DocumentChecksum                    = types.DocumentChecksum
	DocumentReceipt                     = types.DocumentReceipt
	MsgShareDocument                    = types.MsgShareDocument
	MsgSendDocumentReceipt              = types.MsgSendDocumentReceipt
	MsgAddSupportedMetadataSchema       = types.MsgAddSupportedMetadataSchema
	MsgAddTrustedMetadataSchemaProposer = types.MsgAddTrustedMetadataSchemaProposer
)
