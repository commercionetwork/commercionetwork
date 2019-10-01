package cli

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
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
		Short:                      "CommercioMINT transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(getCmdOpenCDP(cdc), getCmdCloseCDP(cdc))

	return txCmd
}

func getCmdOpenCDP(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-cdp [deposit-amount] [timestamp]",
		Short: "Deposit the given amount, open a CDP and give to the signer the consideration liquidity amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			user := cliCtx.GetFromAddress()
			depositedAmount, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			cdpRequest := types.CDPRequest{
				Signer:          user,
				DepositedAmount: depositedAmount,
				Timestamp:       args[1],
			}

			msg := types.NewMsgOpenCDP(cdpRequest)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func getCmdCloseCDP(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "deposit-credits [timestamp]",
		Short: "Close the CDP with the given timestamp, withdraw the liquidity amount from signer's wallet and send the " +
			"previous deposited amount to it",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			user := cliCtx.GetFromAddress()

			msg := types.NewMsgCloseCDP(user, args[0])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
