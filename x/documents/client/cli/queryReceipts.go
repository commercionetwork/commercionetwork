package cli

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
)

func CmdSentReceipts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sent-receipts [user-address]",
		Short: "Get all receipts sent by the user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetSentDocumentsReceiptsRequest{
				Address:    addr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.SentDocumentsReceipts(context.Background(), params)
			if err != nil {
				fmt.Printf("could not get any sent receipt for the given address: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "sent-receipts")

	return cmd
}

func CmdReceivedReceipts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "received-receipts [user-address]",
		Short: "Get all receipts received by the user",

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetReceivedDocumentsReceiptsRequest{
				Address:    addr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.ReceivedDocumentsReceipts(context.Background(), params)
			if err != nil {
				fmt.Printf("could not get any received receipt for the given address: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "received-receipts")

	return cmd
}

func CmdDocumentsReceipts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "documents-receipts [documentUUID]",
		Short: "Get all receipts associated with the given document ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			// consider checking uuid validity for args[0]

			params := &types.QueryGetDocumentsReceiptsRequest{
				UUID:       args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.DocumentsReceipts(context.Background(), params)
			if err != nil {
				fmt.Printf("could not get any received receipt for the given document UUID: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "documents-receipts")

	return cmd
}

func CmdDocumentsUUIDReceipts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "documents-uuid-receipts [documentUUID]",
		Short: "Get all uuid receipts associated with the given document ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			params := &types.QueryGetDocumentsUUIDReceiptsRequest{
				UUID:       args[0],
				Pagination: pageReq,
			}

			res, err := queryClient.DocumentsUUIDReceipts(context.Background(), params)
			if err != nil {
				fmt.Printf("could not get any received receipt for the given document UUID: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "documents-uuid-receipts")

	return cmd
}

func CmdShowReceipt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-receipt [receiptUUID]",
		Short: "Get the receipt with the given UUID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			// consider checking uuid validity for args[0]

			params := &types.QueryGetReceiptRequest{
				UUID:       args[0],
			}

			res, err := queryClient.Receipt(context.Background(), params)
			if err != nil {
				fmt.Printf("could not get any receipt with the given UUID: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}