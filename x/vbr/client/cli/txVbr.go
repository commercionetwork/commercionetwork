package cli

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

)

var _ = strconv.Itoa(0)

func CmdIncrementBlockRewardsPool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [coin-denom] [amount]",
		Short: "Increments the block rewards pool's liquidity by the given amount",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			funder := clientCtx.GetFromAddress()
			amout, e := sdk.ParseCoinNormalized(args[1] + args[0])
			if e != nil {
				return e
			}
			argAmount := []sdk.Coin{amout}

			msg := types.NewMsgIncrementBlockRewardsPool(funder.String(), argAmount)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdSetParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-params [epoch_identifier] [earn_rate]",
		Short: "Set the vbr params with epoch identifier(i.e. \"day\" and the vbr earn rate percentage(Dec))", 
		Long: "Example usage:\n commercionetworkd tx vbr set-params day 0.500000000000000000 --from ",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gov := clientCtx.GetFromAddress()
			epochIdentifier := args[0]
			earnRate, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "invalid earnRate (%s)", err)
			}

			msg := types.NewMsgSetParams(gov.String(), epochIdentifier, earnRate)
			if err2 := msg.ValidateBasic(); err2 != nil {
				return err2
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}