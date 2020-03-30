package cli

import (
	"bufio"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
)

const (
	FlagSign                   = "sign"
	FlagSignStorageURI         = "sign-storage-uri"
	FlagSignSignerInstance     = "sign-signer-instance"
	FlagSignVcrID              = "sign-vcr-id"
	FlagSignCertificateProfile = "sign-certificate-profile"
	FlagSignSdnData            = "sign-sdn-data"
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
		getCmdShareDocument(cdc),
		getCmdSendDocumentReceipt(cdc),
	)

	return txCmd
}

func getCmdShareDocument(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "share [recipient] [document-uuid] [document-metadata-uri] " +
			"[metadata-schema-uri] [metadata-schema-version] " +
			"[document-content-uri] " +
			"[checksum-value] [checksum-algorithm] ",
		Short: "Shares the document with the given recipient address",
		Args:  cobra.RangeArgs(5, 8),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()
			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var checksum *types.DocumentChecksum
			var contentURI string
			if len(args) > 5 {
				contentURI = args[5]
				checksum = &types.DocumentChecksum{
					Value:     args[6],
					Algorithm: args[7],
				}
			}

			document := types.Document{
				ContentURI: contentURI,
				UUID:       args[1],
				Metadata: types.DocumentMetadata{
					ContentURI: args[2],
					Schema: &types.DocumentMetadataSchema{
						URI:     args[3],
						Version: args[4],
					},
				},
				Checksum:   checksum,
				Sender:     sender,
				Recipients: ctypes.Addresses{recipient},
			}

			if viper.GetBool(FlagSign) {
				sdnData, err := types.NewSdnDataFromString(viper.GetString(FlagSignSdnData))
				if err != nil {
					return err
				}

				document.DoSign = &types.DocumentDoSign{
					StorageURI:         viper.GetString(FlagSignStorageURI),
					SignerInstance:     viper.GetString(FlagSignSignerInstance),
					VcrID:              viper.GetString(FlagSignVcrID),
					CertificateProfile: viper.GetString(FlagSignCertificateProfile),
					SdnData:            sdnData,
				}
			}

			msg := types.NewMsgShareDocument(document)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	cmd.Flags().Bool(FlagSign, false, "flag that specifies that we want to sign the document")
	cmd.Flags().String(FlagSignStorageURI, "", "flag that specifies the storage URI to sign")
	cmd.Flags().String(FlagSignSignerInstance, "", "the signer instance needed to sign")
	cmd.Flags().String(FlagSignVcrID, "", "the vcr id needed to sign")
	cmd.Flags().String(FlagSignCertificateProfile, "", "the certificate profile needed to sign")
	cmd.Flags().String(FlagSignSdnData, "", "the sdn data needed to sign")

	return cmd
}

func getCmdSendDocumentReceipt(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-receipt [recipient] [tx-hash] [document-uuid] [proof]",
		Short: "Send the document's receipt with the given recipient address",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()
			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			receipt := types.DocumentReceipt{
				Sender:       sender,
				Recipient:    recipient,
				TxHash:       args[1],
				DocumentUUID: args[2],
				UUID:         uuid.NewV4().String(),
			}

			if len(args) == 4 {
				receipt.Proof = args[3]
			}

			msg := types.NewMsgSendDocumentReceipt(receipt)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}
