package cli

import (
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdSetDidDocument() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-DID-document [document_path]",
		Short: "Sets the DID document for the requesting address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			//keybase, err := keyring.New(sdk.KeyringServiceName(), viper.GetString(flags.FlagKeyringBackend), viper.GetString(flags.FlagHome), os.Stdin)

			argsDDOpath, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}

			// read document from path
			ddoData, err := ioutil.ReadFile(argsDDOpath)
			if err != nil {
				return err
			}

			var didDocument types.DidDocument
			json.Unmarshal(ddoData, &didDocument)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// TODO: check signer is the same as ID in DID document

			// Calculate Proof
			/*signature, err := signDidDocument(clientCtx, didDocument, keybase)

			if err != nil {
				return err
			}

			fromAddressPubkey, err := keybase.KeyByAddress(clientCtx.GetFromAddress())
			if err != nil {
				return fmt.Errorf("could not get keybase for address: %w", err)
			}

			verMeth, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, fromAddressPubkey.GetPubKey())
			if err != nil {
				return fmt.Errorf("could not derive address public key: %w", err)
			}

			proof := types.Proof{
				Type:               types.KeyTypeSecp256k12019,
				Created:            time.Now().String(),
				ProofPurpose:       types.ProofPurposeAuthentication,
				Controller:         clientCtx.GetFromAddress().String(),
				VerificationMethod: verMeth,
				SignatureValue:     signature,
			}*/

			//msg := types.NewMsgSetIdentity(types.ContextDidV1, clientCtx.GetFromAddress().String(), didDocument.PubKeys, proof, didDocument.Service)
			msg := types.NewMsgSetDidDocument(types.ContextDidV1, clientCtx.GetFromAddress().String())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

/*
func signDidDocument(cliCtx client.Context, unsignedDoc types.DidDocument, keybase keyring.Keyring) (string, error) {
	cdc := codec.NewLegacyAmino()
	jsonUnsigned, err := cdc.MarshalJSON(unsignedDoc)
	if err != nil {
		return "", fmt.Errorf("error marshaling doc into json")
	}
	_ = jsonUnsigned
	sign, _, err := keybase.SignByAddress(cliCtx.GetFromAddress(), jsonUnsigned)
	if err != nil {
		return "", fmt.Errorf("failed to sign tx")
	}
	return base64.StdEncoding.EncodeToString(sign), nil
}
*/
