package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

	SdnDataCommonName   = "common_name"
	SdnDataSurname      = "surname"
	SdnDataSerialNumber = "serial_number"
	SdnDataGivenName    = "given_name"
	SdnDataOrganization = "organization"
	SdnDataCountry      = "country"
	InputStringSep      = ","
)

func CmdShareDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use: "share [recipient] [document-uuid] [document-metadata-uri] " +
			"[metadata-schema-uri] [metadata-schema-version] " +
			"[document-content-uri] " +
			"[checksum-value] [checksum-algorithm] ",
		Short: "Shares the document with the given recipient address (First 5 arguments are mandatory)",
		Args:  cobra.RangeArgs(5, 8),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientTxContext(cmd)

			if err != nil {
				return err
			}

			if len(args) == 7 {
				return sdkErr.Wrap(sdkErr.ErrUnauthorized, "Unauthorized number of arguments. If you specify [checksum-value] you have to specify [checksum-algorithm] too")
			}

			sender := cliCtx.GetFromAddress()

			var recipient []string
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

			document := types.Document{
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
				sdnData, err := NewSdnDataFromString(viper.GetString(FlagSignSdnData))
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

/*********************************************************************/
// SdnData represents the SdnData value inside a DocumentDoSign struct.
type SdnData []string

var validSdnData = map[string]struct{}{
	SdnDataCommonName:   {},
	SdnDataSurname:      {},
	SdnDataSerialNumber: {},
	SdnDataGivenName:    {},
	SdnDataOrganization: {},
	SdnDataCountry:      {},
}

// NewSdnDataFromString generates a SdnData slice based on the input string.
// The input string expects a comma-separated value as:
// "common_name,surname,serial_number"
// If empty string is provided, a SdnData with default value will be provided. Default : "serial_number".
func NewSdnDataFromString(input string) (SdnData, error) {
	if input == "" {
		return SdnData{SdnDataSerialNumber}, nil
	}

	var split SdnData = strings.Split(input, InputStringSep)
	err := split.Validate()
	if err != nil {
		return SdnData{}, err
	}

	return split, nil
}

// Validate checks that the SdnData is valid, only accepts value included in
// validSdnData.
func (s SdnData) Validate() error {
	for _, val := range s {
		if _, ok := validSdnData[val]; !ok {
			return fmt.Errorf("sdn_data value \"%s\" is not supported", val)
		}
	}

	return nil
}
