package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	errors "cosmossdk.io/errors"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	uuid "github.com/satori/go.uuid"
)

const (
	FlagSign                   = "sign"
	FlagSignStorageURI         = "sign-storage-uri"
	FlagSignSignerInstance     = "sign-signer-instance"
	FlagSignVcrID              = "sign-vcr-id"
	FlagSignCertificateProfile = "sign-certificate-profile"
	FlagSignSdnData            = "sign-sdn-data"
	FlagFile				   = "file-path"
)

func CmdShareDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use: "share [recipient] [document-uuid] [document-metadata-uri] " +
			"[metadata-schema-uri] [metadata-schema-version] " +
			"[document-content-uri] " +
			"[checksum-value] [checksum-algorithm] ",
		Short: "Shares the document with the given recipient address (First 5 arguments are mandatory)",
		Args:  cobra.RangeArgs(0, 8), // Allow no arguments if file flag is used
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}

			filePath, _ := cmd.Flags().GetString(FlagFile)
			var document types.Document

			if filePath != "" {
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					return err
				}

				// Unmarshal the JSON content
				err = json.Unmarshal(fileContent, &document)
				if err != nil {
					return fmt.Errorf("failed to unmarshal JSON: %w", err)
				}
			} else {
				if len(args) < 5 {
					return errors.Wrap(sdkErr.ErrInvalidRequest, "Unauthorized number of arguments. The first 5 arguments are mandatory if you don't provide a file")
				}

				// this check could be improved by reading args[6] and args[7]
				if len(args) == 7 {
					return errors.Wrap(sdkErr.ErrUnauthorized, "Unauthorized number of arguments. If you specify [checksum-value] you have to specify [checksum-algorithm] too")
				}

				sender := cliCtx.GetFromAddress()

				var recipient []string
				// accepting only one recipient
				recipient = append(recipient, args[0])

				var checksum *types.DocumentChecksum
				var contentURI string
				if len(args) > 5 {
					contentURI = args[5]
					if len(args) > 6 {
						checksum = &types.DocumentChecksum{
							Value:     args[6],
							Algorithm: args[7],
						}
					}
				}

				document = types.Document{
					ContentURI: contentURI,
					UUID:       args[1],
					Metadata: &types.DocumentMetadata{
						ContentURI: args[2],
						Schema: &types.DocumentMetadataSchema{
							URI:     args[3],
							Version: args[4],
						},
					},
					Checksum:   checksum,
					Sender:     sender.String(),
					Recipients: recipient,
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
			}

			msg := types.NewMsgShareDocument(document)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().Bool(FlagSign, false, "flag that specifies that we want to sign the document")
	cmd.Flags().String(FlagSignStorageURI, "", "flag that specifies the storage URI to sign")
	cmd.Flags().String(FlagSignSignerInstance, "", "the signer instance needed to sign")
	cmd.Flags().String(FlagSignVcrID, "", "the vcr id needed to sign")
	cmd.Flags().String(FlagSignCertificateProfile, "", "the certificate profile needed to sign")
	cmd.Flags().String(FlagSignSdnData, "", "the sdn data needed to sign")
	cmd.Flags().String(FlagFile, "", "Path to a JSON file containing the Document")

	return cmd
}

func CmdSendDocumentReceipt() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-receipt [recipient] [tx-hash] [document-uuid] [proof]",
		Short: "Send the document's receipt with the given recipient address ([proof] is optional)",
		Args:  cobra.RangeArgs(3, 4),
		RunE: func(cmd *cobra.Command, args []string) error {

			cliCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sender := cliCtx.GetFromAddress()
			recipient := args[0]
			txHash := args[1]
			documentUUID := args[2]
			UUID := uuid.NewV4().String()

			// empty proof is not valid! consider removing this check and accept all 4 arguments below
			proof := ""
			if len(args) == 4 {
				proof = args[3]
			}

			msg := types.NewMsgSendDocumentReceipt(UUID, sender.String(), recipient, txHash, documentUUID, proof)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
