package cli

import (
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "CommercioTXREWARDS transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(GetCmdIncrementBlockRewardsPool(cdc))

	return txCmd
}

func GetCmdIncrementBlockRewardsPool(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "increment-rewards-pool [coin-denom] [amount]",
		Short: "Increments the block rewards pool's liquidity by the given amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			funder := cliCtx.GetFromAddress()
			coin, _ := sdk.NewIntFromString(args[1])
			amount := sdk.Coin{Denom: args[0], Amount: coin}

			msg := types.NewMsgIncrementBlockRewardsPool(funder, amount)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd = client.PostCommands(cmd)[0]

	return cmd
}
