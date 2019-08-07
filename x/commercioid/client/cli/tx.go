package cli

import (
	types2 "github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

/**
This files contains the functions that take a CLI command and emit a message that will later be held in order to
perform a transaction on the blockchain.
*/

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CommercioID transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(
		GetCmdSetIdentity(cdc),
		GetCmdCreateConnection(cdc),
	)

	return txCmd
}

// GetCmdSetIdentity is the CLI command for sending a SetIdentity transaction
func GetCmdSetIdentity(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-identity [did] [ddo-reference]",
		Short: "Edit an existing identity or add a new one",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			//todo check if we must need NewCLIContextWithFrom
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			account := cliCtx.GetFromAddress()

			msg := types.NewMsgSetIdentity(types2.Did(args[0]), args[1], account)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}

// GetCmdCreateConnection is the CLI command for sending a CreateConnection transaction
func GetCmdCreateConnection(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "create-connection [first-did] [second-did]",
		Short: "Creates a connection between the first and second DIDs",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			account := cliCtx.GetFromAddress()

			msg := types.NewMsgCreateConnection(types2.Did(args[0]), types2.Did(args[1]), account)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
