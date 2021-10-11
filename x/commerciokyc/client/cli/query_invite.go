package cli

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
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

	return cmd
}

func getInvitesFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)
	params := &types.QueryInvitesRequest{}

	res, err := queryClient.Invites(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}
