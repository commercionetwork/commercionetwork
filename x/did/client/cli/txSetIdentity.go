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

func CmdSetIdentity() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-identity [did_document_proposal_path]",
		Short: "Sets the DID document for the requesting address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			argsDDOpath, err := cast.ToStringE(args[0])
			if err != nil {
				return err
			}

			// read DID document proposal from path
			ddoData, err := ioutil.ReadFile(argsDDOpath)
			if err != nil {
				return err
			}

			var didDocumentProposal types.MsgSetDidDocument
			json.Unmarshal(ddoData, &didDocumentProposal)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			if err := didDocumentProposal.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &didDocumentProposal)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
