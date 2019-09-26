package docs

import (
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UserDocumentsData struct {
	User              sdk.AccAddress         `json:"user"`
	SentDocuments     types.Documents        `json:"sent_documents"`
	ReceivedDocuments types.Documents        `json:"received_documents"`
	SentReceipts      types.DocumentReceipts `json:"sent_receipts"`
	ReceivedReceipts  types.DocumentReceipts `json:"received_receipts"`
}

// GenesisState - docs genesis state
type GenesisState struct {
	UsersData                      []UserDocumentsData   `json:"users_data"`
	SupportedMetadataSchemes       types.MetadataSchemes `json:"supported_metadata_schemes"`
	TrustedMetadataSchemaProposers []sdk.AccAddress      `json:"trusted_metadata_schema_proposers"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, data := range data.UsersData {
		keeper.SetUserDocuments(ctx, data.User, data.SentDocuments, data.ReceivedDocuments)
		keeper.SetUserReceipts(ctx, data.User, data.SentReceipts, data.ReceivedReceipts)
	}
	for _, schema := range data.SupportedMetadataSchemes {
		keeper.AddSupportedMetadataScheme(ctx, schema)
	}
	for _, proposer := range data.TrustedMetadataSchemaProposers {
		keeper.AddTrustedSchemaProposer(ctx, proposer)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	users := keeper.GetUsersSet(ctx)

	var usersData []UserDocumentsData
	for _, user := range users {
		sentDocs, err := keeper.GetUserSentDocuments(ctx, user)
		if err != nil {
			panic(err)
		}

		receivedDocs, err := keeper.GetUserReceivedDocuments(ctx, user)
		if err != nil {
			panic(err)
		}

		userData := UserDocumentsData{
			User:              user,
			SentDocuments:     sentDocs,
			ReceivedDocuments: receivedDocs,
			SentReceipts:      keeper.GetUserSentReceipts(ctx, user),
			ReceivedReceipts:  keeper.GetUserReceivedReceipts(ctx, user),
		}
		usersData = append(usersData, userData)
	}

	return GenesisState{
		UsersData:                      usersData,
		SupportedMetadataSchemes:       keeper.GetSupportedMetadataSchemes(ctx),
		TrustedMetadataSchemaProposers: keeper.GetTrustedSchemaProposers(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
