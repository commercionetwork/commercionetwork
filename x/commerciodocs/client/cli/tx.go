package cli

import (
	types2 "github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
		GetCmdStoreDocument(cdc),
		GetCmdShareDocument(cdc),
	)

	return txCmd
}

// GetCmdSetIdentity is the CLI command for sending a StoreDocument transaction
func GetCmdStoreDocument(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "store [identity] [document-reference] [metadata-reference]",
		Short: "Edit an existing identity or add a new one",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			account := cliCtx.GetFromAddress()

			msg := types.NewMsgStoreDocument(account, args[0], args[1], args[2])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
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
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			account := cliCtx.GetFromAddress()

			msg := types.NewMsgShareDocument(account, args[0], args[1], args[2])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
