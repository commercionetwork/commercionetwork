package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/id/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CommercioDOCS id subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		getSetIdentityCommand(cdc),
	)

	return txCmd
}

func getSetIdentityCommand(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "setidentity",
		Short: "sets the identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
}
