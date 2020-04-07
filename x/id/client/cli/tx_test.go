package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"

	thelper "github.com/cosmos/cosmos-sdk/tests"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

func TestGetTxCmd_SetIdentityCommand(t *testing.T) {
	cdc := amino.NewCodec()

	tests := []struct {
		name   string
		errStr string
		flags  []string
	}{
		{
			name:   "privRsaVerKey flag not set",
			errStr: "required flag(s) \"privRsaVerKey\" not set",
			flags: []string{
				fmt.Sprintf("--%s=%s", flagPrivRsaSignKey, "someValue"),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "someValue"),
			},
		},
		{
			name:   "privRsaSignKey flag not set",
			errStr: "required flag(s) \"privRsaSignKey\" not set",
			flags: []string{
				fmt.Sprintf("--%s=%s", flagPrivRsaVerKey, "someValue"),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, "someValue"),
			},
		},
		{
			name:   "from flag not set",
			errStr: "required flag(s) \"from\" not set",
			flags: []string{
				fmt.Sprintf("--%s=%s", flagPrivRsaVerKey, "someValue"),
				fmt.Sprintf("--%s=%s", flagPrivRsaSignKey, "someValue"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			cmd := getSetIdentityCommand(cdc)
			_, out, _ := thelper.ApplyMockIO(cmd)

			cmd.SetArgs(tt.flags)

			cmd.Execute()

			require.True(t, strings.Contains(out.String(), tt.errStr))
		})
	}

}
