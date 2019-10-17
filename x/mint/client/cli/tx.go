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
	txCmd.AddCommand(getCmdOpenCdp(cdc), getCmdCloseCdp(cdc))

	return txCmd
}

func getCmdOpenCdp(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open-cdp [deposit-amount] [timestamp]",
		Short: "Deposit the given amount, open a Cdp and give to the signer the consideration liquidity amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			user := cliCtx.GetFromAddress()
			depositedAmount, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			cdpRequest := types.CdpRequest{
				Signer:          user,
				DepositedAmount: depositedAmount,
				Timestamp:       args[1],
			}

			msg := types.NewMsgOpenCdp(cdpRequest)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func getCmdCloseCdp(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "close-cdp [timestamp]",
		Short: "Close the Cdp with the given timestamp, withdrawing the liquidity amount from signer's wallet and" +
			" sending the previously deposited amount back into it",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			user := cliCtx.GetFromAddress()

			msg := types.NewMsgCloseCdp(user, args[0])
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
