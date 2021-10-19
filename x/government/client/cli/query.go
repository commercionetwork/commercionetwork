package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group id queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		CmdGetGovernmentAddr(),
	)
	return cmd
}

func CmdGetGovernmentAddr() *cobra.Command {
	return &cobra.Command{
		Use:   "gov-address",
		Short: "Get the government address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryGovernmentAddress)
			res, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("could not get government address: %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
