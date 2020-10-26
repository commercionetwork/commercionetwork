package cli

import (
	"bufio"
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		scheduleUpgradeHeightCmd(cdc),
		scheduleUpgradeTimeCmd(cdc),
		deleteUpgradeCmd(cdc),
	)

	return txCmd
}

func scheduleUpgradeHeightCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-height [NAME] [HEIGHT] [INFO]",
		Short: "Schedule an upgrade with block height",
		Long:  "Example usage:\n cncli tx upgrade schedule-height testUpgrade 10 this_is_just_a_test --from ...",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scheduleUpgradeCmdFunc(cmd, args, cdc, true)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

const (
	// TimeFormat specifies ISO UTC format for submitting the time for a new upgrade proposal
	TimeFormat = "2006-01-02T15:04:05Z"
)

func scheduleUpgradeTimeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-time [NAME] [TIME] [INFO]",
		Short: fmt.Sprint("Schedule an upgrade with time format ", TimeFormat),
		Long:  "Example usage:\n cncli tx upgrade schedule-time testUpgrade 2020-10-23T15:21:05Z this_is_just_a_test --from ",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scheduleUpgradeCmdFunc(cmd, args, cdc, false)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func scheduleUpgradeCmdFunc(cmd *cobra.Command, args []string, cdc *codec.Codec, useHeight bool) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	proposer := cliCtx.GetFromAddress()

	name := args[0]

	var height int64
	var upgradeTime time.Time
	var err error
	if useHeight {
		height, err = strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return err
		}
	} else {
		upgradeTime, err = time.Parse(TimeFormat, args[1])
		if err != nil {
			return err
		}
	}

	info := args[2]

	plan := upgrade.Plan{
		Name:   name,
		Height: height,
		Time:   upgradeTime,
		Info:   info,
	}

	msg := types.MsgScheduleUpgrade{
		Proposer: proposer,
		Plan:     plan,
	}

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}

func deleteUpgradeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete any scheduled upgrade",
		Long:  "Example usage:\n cncli tx upgrade delete --from ...",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteUpgradeCmdFunc(cmd, args, cdc)
		},
	}

	cmd = flags.PostCommands(cmd)[0]

	return cmd
}

func deleteUpgradeCmdFunc(cmd *cobra.Command, _ []string, cdc *codec.Codec) error {
	inBuf := bufio.NewReader(cmd.InOrStdin())
	cliCtx := context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

	proposer := cliCtx.GetFromAddress()

	msg := types.MsgDeleteUpgrade{
		Proposer: proposer,
	}

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
}
