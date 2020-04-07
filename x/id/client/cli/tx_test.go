package cli

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

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
			},
		},
		{
			name:   "privRsaSignKey flag not set",
			errStr: "required flag(s) \"privRsaSignKey\" not set",
			flags: []string{
				fmt.Sprintf("--%s=%s", flagPrivRsaVerKey, "someValue"),
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

func TestPrivKey(t *testing.T) {
	file := "/Users/jgimeno/.ssh/test"

	priv, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	decode, _ := pem.Decode(priv)

	fmt.Printf("%v", decode)
}
