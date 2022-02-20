package cli

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdMemberships() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "memberships",
		Short: "Get all memberships",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMembershipsFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "memberships")

	return cmd
}

func getMembershipsFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)
	pageReq, err := client.ReadPageRequest(cmd.Flags())
	if err != nil {
		return err
	}

	params := &types.QueryMembershipsRequest{Pagination: pageReq}

	res, err := queryClient.Memberships(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}

func CmdMembership() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "membership [user]",
		Short: "Get user membership",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getMembershipFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getMembershipFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)
	// TODO check valid address
	params := &types.QueryMembershipRequest{
		Address: args[0],
	}

	res, err := queryClient.Membership(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}
