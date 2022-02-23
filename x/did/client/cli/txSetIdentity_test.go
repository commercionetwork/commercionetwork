package cli_test

/*
func TestSetIdentity(t *testing.T) {
	net := network.New(t)
	val := net.Validators[0]
	ctx := val.ClientCtx

	ddoDoc := `{
		  "@context": "https://www.w3.org/ns/did/v1",
		  "id": "` + val.Address.String() + `",
		  "publicKey": [
			{
			  "id": "` + val.Address.String() + `#keys-1",
			  "type": "RsaVerificationKey2018",
			  "controller": "` + val.Address.String() + `",
			  "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+QNCVuyIF4l+n2fYT99f\nv6qPr3/Xe8cD8xoaQrWthg5VMk6s0WhbyZMg7JD0i6I9NTCt0IOoz0N7N61e+igb\nK0wnsHCE+/ZZkcQ14pU4NneQPXM/z8MBOKXxfoFbGEeAPDf1bUmqrYQ7z4WuEpJo\ngTfojC8EpE0pufrrZkVUXKos2hhTEAJvGXkpa+TRwiaCEC8q8KjKaRYuMEVJe5yN\nNV5pkDiz6hAwEuE3FCy4y2h1FqPIfZNQF5LPdpZ2fXq19O1wx0S+XOxf3KXzX4b6\n2BtubNalDKigXTrHK2RQlw8z83dnoX8Vwek2vXoz0P6rGeGxnDYdECgtofBnd4eL\n0QIDAQAB\n-----END PUBLIC KEY-----\n"
			},
			{
			  "id": "` + val.Address.String() + `#keys-2",
			  "type": "RsaSignatureKey2018",
			  "controller": "` + val.Address.String() + `",
			  "publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAwGvTpscbNrJtmh9AwUZZ\ndbNgMs1Y4H8CUq4eK0Ddxzr6xGQbvq9FI8tpkUVOViR9p9OhPKdp0WR+EWHVGhxf\n1uIE42NO1y4d/IE2WjA6/TiX897enLWsuTqOUt9Z7FKtcVtfANhn5miZQlhanO6h\n4Cq4mOF/KRrPAaFt2ZU2M56QV+9ZtJ0uzPPd5p8yBIGi861EIHgbMXLoYKoUVGt2\nBNJcVNVFUklwwVqGvXxa1VvkPIMiXkGUN28JyQhrN+f0HaqDOogDYFDjS/d9X2D0\n/0XO0gCDzQyWTUCaxiP8l9dE9QcFlkhBDcoCt1TsBNi/iG6ZiPeiSJvPsivwg9XS\n8qAWUSKN0xn0yLKKji7ipa04lbBl+bEMg5u6vhzmFuYADcUI7ov9FCu8LGBe4ybX\nPsw3vbMtPrOetzremZDYYUuE+PJAcseNPKJ1xXMC4Sl6cPrA8ZIsJ+W8BMQ7zbDO\nWajP50p6281qg74/ftyB/Gt1L5kFRnJV/z+OAqrpAHWHdrSSTQvGKSiABT1xLxO9\nTYJnfVAxDeTY2trDnJa7w/yBDGxjvILp0ib5PmEXxwhX2EUYyotokzLltVG9nlgf\nA2BdYFefKOiEFJ6+LIZORy1i83PSAxrbU8+AneM2IUH851cSGV4jIVuJJfbA7lF5\nJb2BB1sjnrGYiiByN4IXdLECAwEAAQ==\n-----END PUBLIC KEY-----\n"
			}
		  ]
		}
	`

	ddoFile := tutil.WriteToNewTempFile(t, ddoDoc)
	fields := []string{ddoFile.Name()}

	for _, tc := range []struct {
		desc string
		args []string
		err  error
		code uint32
	}{
		{
			desc: "valid",
			args: []string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, val.Address.String()),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(net.Config.BondDenom, sdk.NewInt(10))).String()),
			},
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			args := []string{}
			args = append(args, fields...)
			args = append(args, tc.args...)
			out, err := clitestutil.ExecTestCLICmd(ctx, cli.CmdSetIdentity(), args)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				var resp sdk.TxResponse
				require.NoError(t, ctx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &resp))
				require.Equal(t, tc.code, resp.Code)
			}
		})
	}
}
*/
