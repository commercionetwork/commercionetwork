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
		Short: "Sets a JSON-LD DID document for the requesting address, reading it from the path specified as first parameter. The file must conform to the rules of Decentralized Identitfiers (DIDs) v1.0 plus additional rules defined by commercionetwork. Please refer to https://www.w3.org/TR/2021/PR-did-core-20210803/ and https://docs.commercio.network/x/did/",
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

			var didDocumentProposal types.DidDocument
			json.Unmarshal(ddoData, &didDocumentProposal)

			if err := didDocumentProposal.Validate(); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msgSetIdentity := types.MsgSetIdentity{
				DidDocument: &didDocumentProposal,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msgSetIdentity)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
