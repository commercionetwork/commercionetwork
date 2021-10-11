package cli

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetConversionRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "conversion-rate",
		Short: "Display the current conversion rate",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryConversionRate{}

			res, err := queryClient.ConversionRate(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetFreezePeriod() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "freeze-period",
		Short: "Display the current freeze period",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryFreezePeriod{}

			res, err := queryClient.FreezePeriod(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
