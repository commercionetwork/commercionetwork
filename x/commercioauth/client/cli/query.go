package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetCmdReadAccount(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "account [address]",
		Short: "Retrieves the account information of the account with the given address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			name := args[0]

			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/account/%s", queryRoute, name), nil)
			if err != nil {
				fmt.Printf("Could not get account for address %s: \n %s", string(name), err)
				return nil
			}

			fmt.Println(string(res))

			return nil
		},
	}
}
