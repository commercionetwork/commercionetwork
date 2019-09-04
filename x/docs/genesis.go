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
	UsersData []UserDocumentsData `json:"users_data"`
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
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	users, err := keeper.GetUsersSet(ctx)
	if err != nil {
		panic(err)
	}

	var usersData []UserDocumentsData
	for _, user := range users {
		userData := UserDocumentsData{
			User:              user,
			SentDocuments:     keeper.GetUserSentDocuments(ctx, user),
			ReceivedDocuments: keeper.GetUserReceivedDocuments(ctx, user),
			SentReceipts:      keeper.GetUserSentReceipts(ctx, user),
			ReceivedReceipts:  keeper.GetUserReceivedReceipts(ctx, user),
		}
		usersData = append(usersData, userData)
	}

	return GenesisState{
		UsersData: usersData,
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(_ GenesisState) error {
	return nil
}
