package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

func CmdGetGovernmentAddr() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov-address",
		Short: "Get the government address",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryGovernmentAddrRequest{}

			res, err := queryClient.GovernmentAddr(context.Background(), params)

			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
