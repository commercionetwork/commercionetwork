package cli

import (
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
		GetCmdShareDocument(cdc),
	)

	return txCmd
}

// GetCmdShareDocument is the CLI command for sending a ShareDocument transaction
func GetCmdShareDocument(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "share [recipient] [document-uuid] [document-content-uri] [metadata-content-uri] [metadata-schema-uri] [schema-version] " +
			"[computation-proof] [checksum-value] [checksum-algorithm]",
		Short: "Shares the document with the given recipient address",
		Args:  cobra.ExactArgs(9),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			sender := cliCtx.GetFromAddress()
			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgShareDocument(sender, recipient, args[1], args[2], args[3], args[4], args[5], args[6],
				args[7], args[8])
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
