package cli

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

/*
getCmdTspMemberships(cdc),
*/

func CmdTrustedServiceProviders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trusted-service-providers",
		Short: "Get all trusted services providers",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getTrustedServiceProvidersFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getTrustedServiceProvidersFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)
	params := &types.QueryTspsRequest{}

	res, err := queryClient.Tsps(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}

func CmdPoolFunds() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pool-funds",
		Short: "Get the pool funds amounts",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPoolFundsFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getPoolFundsFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)
	params := &types.QueryFundsRequest{}

	res, err := queryClient.Funds(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}
