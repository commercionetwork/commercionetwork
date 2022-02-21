package cli

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdShowHistory() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-history [id]",
		Short: "Shows the list of DID document updates for the specified id, in chronologial order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			id := string(args[0])

			params := &types.QueryResolveIdentityHistoryRequest{
				ID: id,
			}

			res, err := queryClient.IdentityHistory(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintObjectLegacy(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
