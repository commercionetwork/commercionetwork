package cli

import (
	"commercio-network/types"
	"commercio-network/x/commerciodocs"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
)

// GetCmdSetIdentity is the CLI command for sending a StoreDocument transaction
func GetCmdStoreDocument(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "store [identity] [document-reference] [metadata-reference]",
		Short: "Edit an existing identity or add a new one",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			//txBldr := authtxb.NewTxBuilderFromCLI()
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			msg := commerciodocs.NewMsgStoreDocument(account, types.Did(args[0]), args[1], args[2])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

// GetCmdShareDocument is the CLI command for sending a ShareDocument transaction
func GetCmdShareDocument(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "share [document-reference] [sender-identity] [receiver-identity]",
		Short: "Shares the document with the given reference between the sender identity and the receiver identity",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			//txBldr := authtxb.NewTxBuilderFromCLI()
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			msg := commerciodocs.NewMsgShareDocument(account, args[0], types.Did(args[1]), types.Did(args[2]))
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
