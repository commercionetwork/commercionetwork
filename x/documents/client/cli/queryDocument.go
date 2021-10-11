package cli

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

// ----------------------------------
// --- Documents
// ----------------------------------
/*
func CmdListDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-document",
		Short: "list all document",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllDocumentRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.DocumentAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
*/
func CmdShowDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-document [doc-uuid]",
		Short: "shows a document",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			id := string(args[0])

			params := &types.QueryGetDocumentRequest{
				UUID: id,
			}

			res, err := queryClient.Document(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdSentDocuments() *cobra.Command {
	return &cobra.Command{
		Use:   "sent-documents [user-address]",
		Short: "Get all documents sent by user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentDocuments, args[0])
			res, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get sent documents by user: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

func CmdReceivedDocuments() *cobra.Command {
	return &cobra.Command{
		Use:   "received-documents [user-address]",
		Short: "Get all documents received by user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryReceivedDocuments, args[0])
			res, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				fmt.Printf("Could not get received documents by user: \n %s", err)
			}

			fmt.Println(string(res))

			return nil
		},
	}
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func CmdSentReceipts() *cobra.Command {
	return &cobra.Command{
		Use:   "sent-receipts [user-address]",
		Short: "Get all the receipts sent from the specified user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentReceipts, args[0])
			res, _, err2 := clientCtx.QueryWithData(route, nil)
			if err2 != nil {
				fmt.Printf("Could not get any sent receipt for the given user: \n %s", err2)
			}

			fmt.Print(string(res))
			return nil
		},
	}
}

func CmdReceivedReceipts() *cobra.Command {
	return &cobra.Command{
		Use:   "received-receipts [user-address] [[doc-uuid]]",
		Short: "Get the document receipt associated with given document uuid",
		Args:  cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			addr, uuid := args[0], ""
			if len(args) == 2 {
				uuid = args[1]
			}

			route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryReceivedReceipts, addr, uuid)
			res, _, err2 := clientCtx.QueryWithData(route, nil)
			if err2 != nil {
				fmt.Printf("Could not get any received receipt for the given user or uuid: \n %s", err2)
			}

			fmt.Print(string(res))
			return nil
		},
	}
}
