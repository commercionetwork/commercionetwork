package cli

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		getCmdAccrediter(cdc),
		getCmdSigners(cdc),
	)

	return cmd
}

func getCmdAccrediter(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "accrediter [user-address]",
		Short: "Get the accrediter of user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetAccrediter, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get sent documents by user: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func getCmdSigners(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "signers",
		Short: "Get all the trustworthy signers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryGetSigners, args[0])
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get trustworthy signers: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
