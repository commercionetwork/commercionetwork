package cli

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetEtps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-etps [user-addr]",
		Short: "Get all opened ETPs for an user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getEtpsFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getEtpsFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)

	sender, err := sdk.AccAddressFromBech32(args[0])
	if err != nil {
		return err
	}

	pageReq, err := client.ReadPageRequest(cmd.Flags())
	if err != nil {
		return err
	}

	params := &types.QueryEtpRequestByOwner{
		Owner: sender.String(),
		Pagination: pageReq,
	}

	res, err := queryClient.EtpByOwner(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}

func CmdGetAllEtps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-all-etps",
		Short: "Get all opened ETPs for an user",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryEtpsRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.Etps(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdGetEtp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-etp [id]",
		Short: "Get opened ETP by an id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getEtpFunc(cmd, args)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func getEtpFunc(cmd *cobra.Command, args []string) error {
	clientCtx := client.GetClientContextFromCmd(cmd)

	queryClient := types.NewQueryClient(clientCtx)

	params := &types.QueryEtpRequest{
		ID: string(args[0]),
	}

	res, err := queryClient.Etp(context.Background(), params)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)

}