package cli

import (
	"bufio"
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Pricefeed transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(GetCmdSetPrice(cdc), GetCmdAddOracle(cdc), GetCmdBlacklistDenom(cdc))

	return txCmd
}

func GetCmdSetPrice(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "set-price [token-name] [token-price] [expiry]",
		Short: "set price for a given token",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			tokenPrice, err := sdk.NewDecFromStr(args[1])
			if err != nil {
				return sdkErr.Wrap(sdkErr.ErrInvalidRequest, (err.Error()))
			}

			expiry, ok := sdk.NewIntFromString(args[2])
			if !ok {
				return sdkErr.Wrap(sdkErr.ErrInvalidRequest, (fmt.Sprintf("Invalid expiration height, %s", args[2])))
			}

			price := types.NewPrice(args[0], tokenPrice, expiry)

			oracle := cliCtx.GetFromAddress()
			msg := types.NewMsgSetPrice(price, oracle)

			err2 := msg.ValidateBasic()
			if err2 != nil {
				return err2
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdAddOracle cli command for posting prices.
func GetCmdAddOracle(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "add-oracle [oracle-address]",
		Short: "add a trusted oracle to the oracles' list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			oracle, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			signer := cliCtx.GetFromAddress()
			msg := types.NewMsgAddOracle(signer, oracle)

			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// GetCmdBlacklistDenom cli command for blacklisting denoms.
func GetCmdBlacklistDenom(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "blacklist-denom [denom]",
		Short: "blacklists a denom, prevent oracle from setting price for it",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getCmdBlacklistDenomFunc(cmd, cdc, args)
		},
	}
}

func getCmdBlacklistDenomFunc(cmd *cobra.Command, cdc *codec.Codec, args []string) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	denom := args[0]
	signer := cliCtx.GetFromAddress()
	msg := types.NewMsgBlacklistDenom(signer, denom)

	err := msg.ValidateBasic()
	if err != nil {
		return err
	}

	return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
}
