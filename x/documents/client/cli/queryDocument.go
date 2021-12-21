package cli

import (
	"context"
	//"errors"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
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
	cmd := &cobra.Command{
		Use:   "sent-documents [user-address]",
		Short: "Get all documents sent by user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			addr, e := sdk.AccAddressFromBech32(args[0])
			if e != nil {
				return e
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			//pageReq.Limit = uint64(10)
			params := &types.QueryGetSentDocumentsRequest{
				Address: addr.String(),
				Pagination: pageReq,
			}
			
			//route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentDocuments, args[0])
			res, err := queryClient.SentDocuments(context.Background(), params)
			if err != nil {
				return sdkErr.Wrap(sdkErr.ErrLogic, fmt.Sprintf("Could not get sent documents by user: \n %s", err))
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdReceivedDocuments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "received-documents [user-address]",
		Short: "Get all documents received by user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			addr, e := sdk.AccAddressFromBech32(args[0])
			if e != nil {
				return e
			}
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			//pageReq.Limit = uint64(10)
			params := &types.QueryGetReceivedDocumentRequest{
				Address: addr.String(),
				Pagination: pageReq,
			}

			//route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryReceivedDocuments, args[0])
			res, err := queryClient.ReceivedDocument(context.Background(), params)
			if err != nil {
				fmt.Printf("Could not get received documents by user: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func CmdSentReceipts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sent-receipts [user-address]",
		Short: "Get all the receipts sent from the specified user",
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
			//pageReq.Limit = uint64(10)
			params := &types.QueryGetSentDocumentsReceiptsRequest{
				Address: addr.String(),
				Pagination: pageReq,
			}

			//route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QuerySentReceipts, args[0])
			res, err := queryClient.SentDocumentsReceipts(context.Background(), params)
			if err != nil {
				fmt.Printf("Could not get any sent receipt for the given user: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdReceivedReceipts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "received-receipts [user-address]",
		Short: "Get the document receipt associated with given address",
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
			//pageReq.Limit = uint64(10)
			params := &types.QueryGetReceivedDocumentsReceiptsRequest{
				Address: addr.String(),
				Pagination: pageReq,
			}

			//route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, types.QueryReceivedReceipts, addr, uuid)
			res, err := queryClient.ReceivedDocumentsReceipts(context.Background(), params)
			if err != nil {
				fmt.Printf("Could not get any received receipt for the given user or uuid: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
