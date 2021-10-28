package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var _ = strconv.Itoa(0)

const (
	// TimeFormat specifies ISO UTC format for submitting the time for a new upgrade proposal
	TimeFormat = "2006-01-02T15:04:05Z"
)

func scheduleUpgradeCmdFunc(cmd *cobra.Command, args []string, useHeight bool) error {
	cliCtx, e := client.GetClientTxContext(cmd)
	if e != nil {
		return e
	}

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
		upgradeTime = upgradeTime.UTC()
		if err != nil {
			return err
		}
	}

	info := args[2]
	plan := upgradetypes.Plan{
		Name:   name,
		Height: height,
		Time:   upgradeTime,
		Info:   info,
	}
	
	msg := types.NewMsgScheduleUpgrade(proposer.String(), plan)
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}

func CmdScheduleUpgradeTime() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-time [NAME] [TIME] [INFO]",
		Short: fmt.Sprint("Schedule an upgrade with UTC time using format ", TimeFormat),
		Long:  "Example usage:\n commercionetworkd tx upgrade schedule-time testUpgrade 2020-10-23T15:21:05Z this_is_just_a_test --from ",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scheduleUpgradeCmdFunc(cmd, args, false)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdScheduleUpgradeHeight() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "schedule-height [NAME] [HEIGHT] [INFO]",
		Short: "Schedule an upgrade with block height",
		Long:  "Example usage:\n commercionetworkd tx upgrade schedule-height testUpgrade 10 this_is_just_a_test --from ...",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			return scheduleUpgradeCmdFunc(cmd, args, true)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func CmdDeleteUpgrade() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete any scheduled upgrade",
		Long:  "Example usage:\n commercionetworkd tx upgrade delete --from ...",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteUpgradeCmdFunc(cmd, args)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
    return cmd
}

func deleteUpgradeCmdFunc(cmd *cobra.Command, _ []string) error {
	cliCtx, e := client.GetClientTxContext(cmd)
	if e != nil {
		return e
	}

	proposer := cliCtx.GetFromAddress()
	msg := types.NewMsgDeleteUpgrade(proposer.String())

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), msg)
}
