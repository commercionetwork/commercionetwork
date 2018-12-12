package cli

import (
	"github.com/spf13/cobra"

	"commercio-network/x/commercioid"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

/**
This files contains the functions that take a CLI command and emit a message that will later be held in order to
perform a transaction on the blockchain.
*/

// GetCmdSetIdentity is the CLI command for sending a SetIdentity transaction
func GetCmdSetIdentity(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-identity [did] [ddo-reference]",
		Short: "Edit an existing identity or add a new one",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := commercioid.NewMsgSetIdentity(args[0], args[1], account)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
