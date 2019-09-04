package cli

import (
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CommercioDOCS transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdShareDocument(cdc),
		GetCmdSendDocumentReceipt(cdc),
	)

	return txCmd
}

// GetCmdShareDocument is the CLI command for sending a ShareDocument transaction
func GetCmdShareDocument(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "share [recipient] [document-uuid] [document-metadata-uri] " +
			"[metadata-schema-uri] [metadata-schema-version] [metadata-verification-proof] " +
			"[document-content-uri]" +
			"[checksum-value] [checksum-algorithm]",
		Short: "Shares the document with the given recipient address",
		Args:  cobra.RangeArgs(6, 9),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()
			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var contentUri, checksumValue, checksumAlgorithm string
			if len(args) > 6 {
				contentUri = args[6]
				checksumValue = args[7]
				checksumAlgorithm = args[8]
			}

			document := types.Document{
				Sender:     sender,
				Recipient:  recipient,
				ContentUri: contentUri,
				Uuid:       args[1],
				Metadata: types.DocumentMetadata{
					ContentUri: args[2],
					Schema: types.DocumentMetadataSchema{
						Uri:     args[3],
						Version: args[4],
					},
					Proof: args[5],
				},
				Checksum: types.DocumentChecksum{
					Value:     checksumValue,
					Algorithm: checksumAlgorithm,
				},
			}

			msg := types.NewMsgShareDocument(document)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

func GetCmdSendDocumentReceipt(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-receipt [recipient] [tx-hash] [document-uuid] [proof]",
		Short: "Send the document's receipt with the given recipient address",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()
			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			receipt := types.DocumentReceipt{
				Sender:       sender,
				Recipient:    recipient,
				TxHash:       args[1],
				DocumentUuid: args[2],
				Proof:        args[3],
			}

			msg := types.NewMsgDocumentReceipt(receipt)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
