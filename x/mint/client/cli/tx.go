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
	txCmd.AddCommand(getCmdDepositToken(cdc), getCmdWithdrawToken(cdc))

	return txCmd
}

func getCmdDepositToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-token [amount]",
		Short: "Deposit the given tokens and send the corresponding value in commercio credits to user's wallet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			user := cliCtx.GetFromAddress()
			token, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgDepositToken{
				Signer: user,
				Tokens: token,
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

func getCmdWithdrawToken(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit-credits [amount]",
		Short: "Deposit the given credits and send the corresponding value in commercio tokens to user's wallet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			user := cliCtx.GetFromAddress()
			token, err := sdk.ParseCoins(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgWithdrawToken{
				Signer: user,
				Amount: token,
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
