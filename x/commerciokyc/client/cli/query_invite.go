package cli

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdGetInvites() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invites",
		Short: "Get all membership invitations",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getInvitesFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "invites")

	return cmd
}

func getInvitesFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)

	pageReq, err := client.ReadPageRequest(cmd.Flags())
	if err != nil {
		return err
	}

	params := &types.QueryInvitesRequest{Pagination: pageReq}

	res, err := queryClient.Invites(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}

func CmdGetInvite() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invite [user]",
		Short: "Get user invite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getInviteFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getInviteFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)
	addr, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	params := &types.QueryInviteRequest{
		Address: addr.String(),
	}

	res, err := queryClient.Invite(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}
