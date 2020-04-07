package cli

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/id/types"
)

const flagPrivRsaVerKey = "privRsaVerKey"
const flagPrivRsaSignKey = "privRsaSignKey"

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CommercioDOCS id subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		getSetIdentityCommand(cdc),
	)

	return txCmd
}

func getSetIdentityCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setidentity",
		Short: "sets the identity",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			verPubKey, err := getVerificationPublicKey(cliCtx, viper.GetString(flagPrivRsaVerKey))
			if err != nil {
				return fmt.Errorf("error extracting public from private key verification path: %s", err)
			}

			signPubKey, err := getSignPublicKey(cliCtx, viper.GetString(flagPrivRsaSignKey))
			if err != nil {
				return fmt.Errorf("error extracting public from private key sign path: %s", err)
			}

			unsignedDoc := types.DidDocumentUnsigned{
				Context: types.ContextDidV1,
				ID:      cliCtx.GetFromAddress(),
				PubKeys: types.PubKeys{
					verPubKey,
					signPubKey,
				},
			}

			fmt.Printf("%+v", unsignedDoc)

			return nil
		},
	}

	cmd.Flags().String(flagPrivRsaSignKey, "", "the path of the pem rsa sign key")
	cmd.MarkFlagRequired(flagPrivRsaSignKey)

	cmd.Flags().String(flagPrivRsaVerKey, "", "the path of the pem rsa verification key")
	cmd.MarkFlagRequired(flagPrivRsaVerKey)

	cmd = flags.PostCommands(cmd)[0]
	cmd.MarkFlagRequired(flags.FlagFrom)

	return cmd
}

func getVerificationPublicKey(cliCtx context.CLIContext, path string) (types.PubKey, error) {
	pk, err := types.LoadRSAPrivKeyFromDisk(path)
	if err != nil {
		return types.PubKey{}, err
	}

	fromAddress := cliCtx.GetFromAddress()
	verPubKey := types.PubKey{
		ID:         fmt.Sprintf("%s#keys-1", fromAddress.String()),
		Type:       types.KeyTypeRsaVerification,
		Controller: fromAddress,
		PublicKey:  types.PublicKeyToPemString(&pk.PublicKey),
	}

	return verPubKey, nil
}

func getSignPublicKey(cliCtx context.CLIContext, path string) (types.PubKey, error) {
	pk, err := types.LoadRSAPrivKeyFromDisk(path)
	if err != nil {
		return types.PubKey{}, err
	}

	fromAddress := cliCtx.GetFromAddress()
	verPubKey := types.PubKey{
		ID:         fmt.Sprintf("%s#keys-2", fromAddress.String()),
		Type:       types.KeyTypeRsaSignature,
		Controller: fromAddress,
		PublicKey:  types.PublicKeyToPemString(&pk.PublicKey),
	}

	return verPubKey, nil
}
