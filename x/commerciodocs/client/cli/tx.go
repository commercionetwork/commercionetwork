package cli

import (
	internal "github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
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
	)

	return txCmd
}

// GetCmdShareDocument is the CLI command for sending a ShareDocument transaction
func GetCmdShareDocument(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
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

			document := internal.Document{
				Sender:     sender,
				Recipient:  recipient,
				ContentUri: contentUri,
				Uuid:       args[1],
				Metadata: internal.DocumentMetadata{
					ContentUri: args[2],
					Schema: internal.DocumentMetadataSchema{
						Uri:     args[3],
						Version: args[4],
					},
					Proof: args[5],
				},
				Checksum: internal.DocumentChecksum{
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
}
