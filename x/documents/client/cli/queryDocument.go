package cli

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"
)

func CmdShowDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-document [documentUUID]",
		Short: "Get the document with the given ID",
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
		Short: "Get all documents sent by the user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			// use err name for variable
			addr, e := sdk.AccAddressFromBech32(args[0])
			if e != nil {
				return e
			}

			// missing AddPaginationFlagsToCmd below!
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.QueryGetSentDocumentsRequest{
				Address:    addr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.SentDocuments(context.Background(), params)
			if err != nil {
				return sdkErr.Wrap(sdkErr.ErrLogic, fmt.Sprintf("could not get any sent document for the given address: \n %s", err))
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "sent-documents")

	return cmd
}

func CmdUUIDDocuments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid-documents [user-address]",
		Short: "Get all uuid-documents sent by the user",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			// use err name for variable
			addr, e := sdk.AccAddressFromBech32(args[0])
			if e != nil {
				return e
			}

			// missing AddPaginationFlagsToCmd below!
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}
			params := &types.QueryGetUUIDDocumentsRequest{
				Address:    addr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.UUIDDocuments(context.Background(), params)
			if err != nil {
				return sdkErr.Wrap(sdkErr.ErrLogic, fmt.Sprintf("could not get any sent document for the given address: \n %s", err))
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "uuid-documents")

	return cmd
}

func CmdReceivedDocuments() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "received-documents [user-address]",
		Short: "Get all documents received by the user",
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

			params := &types.QueryGetReceivedDocumentRequest{
				Address:    addr.String(),
				Pagination: pageReq,
			}

			res, err := queryClient.ReceivedDocument(context.Background(), params)
			if err != nil {
				fmt.Printf("could not get any received documents by user: \n %s", err)
			}

			return clientCtx.PrintProto(res)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "received-documents")

	return cmd
}
