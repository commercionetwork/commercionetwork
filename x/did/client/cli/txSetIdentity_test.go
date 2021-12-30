package cli_test

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/commercionetwork/commercionetwork/testutil/network"
// 	"github.com/commercionetwork/commercionetwork/x/did/client/cli"
// 	"github.com/cosmos/cosmos-sdk/client/flags"
// 	tutil "github.com/cosmos/cosmos-sdk/testutil"
// 	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// )

// // TODO
// func TestSetIdentity(t *testing.T) {
// 	net := network.New(t)
// 	val := net.Validators[0]
// 	ctx := val.ClientCtx

// 	ddoDoc := `{
// 		"@context": [
// 			"https://www.w3.org/ns/did/v1",
// 			"https://w3id.org/security/suites/ed25519-2018/v1",
// 			"https://w3id.org/security/suites/x25519-2019/v1"
// 		],
// 		"id": "` + val.Address.String() + `",
// 		"verificationMethod": [
// 			{
// 				"type": "RsaVerificationKey2018",
// 				"id": "` + val.Address.String() + `#keys-1",
// 				"controller": "` + val.Address.String() + `",
// 				"publicKeyMultiBase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
// 			},
// 			{
// 				"type": "RsaSignature2018",
// 				"id": "` + val.Address.String() + `#keys-2",
// 				"controller": "` + val.Address.String() + `",
// 				"publicKeyMultiBase": "mMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
// 			}
// 		],
// 		"authentication": [
// 			"` + val.Address.String() + `#keys-1"
// 		],
// 		"keyAgreement": [
// 			"` + val.Address.String() + `#keys-2"
// 		],
// 		"service": [
// 			{
// 				"id": "A",
// 				"type": "agent",
// 				"serviceEndpoint": "https://commerc.io/agent/serviceEndpoint/"
// 			}
// 		]
// 	}`

// 	ddoFile := tutil.WriteToNewTempFile(t, ddoDoc)
// 	fields := []string{ddoFile.Name()}

// 	for _, tc := range []struct {
// 		desc string
// 		args []string
// 		err  error
// 		code uint32
// 	}{
// 		{
// 			desc: "valid",
// 			args: []string{
// 				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
// 				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
// 				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
// 				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
// 			},
// 		},
// 	} {
// 		tc := tc
// 		t.Run(tc.desc, func(t *testing.T) {
// 			args := []string{}
// 			args = append(args, fields...)
// 			args = append(args, tc.args...)
// 			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSetIdentity(), args)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.NoError(t, err)
// 				var resp sdk.TxResponse
// 				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
// 				require.Equal(t, tc.code, resp.Code)
// 			}
// 		})
// 	}
// }
