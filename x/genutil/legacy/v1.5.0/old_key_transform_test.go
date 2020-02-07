package v1_5_0

import (
	"encoding/base64"
	"testing"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func Test_getPk(t *testing.T) {
	tests := []struct {
		name     string
		pubKey   PubKey
		mustFail bool
	}{
		{
			"random text, no meaning",
			PubKey{
				Value: "stuff",
			},
			true,
		},
		{
			"known good public key",
			PubKey{
				Value: "A01/bky+d3gmuZ+e5RMpRXjx3nFPVXc1om+/Qno6taQ6",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.mustFail {
				require.Panics(t, func() { getPk(tt.pubKey) })
			} else {
				//get pk
				pk := getPk(tt.pubKey)

				// transform it to bech32
				pks, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pk)
				require.NoError(t, err)

				// transform it back to crypto.Pubkey
				orpk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, pks)
				require.NoError(t, err)

				// explicitly handle it as a secp256k1.PubKeySecp256k1
				secpk := orpk.(secp256k1.PubKeySecp256k1)

				newpk := make([]byte, len(secpk))
				copy(newpk, secpk[:])

				require.Equal(t, base64.StdEncoding.EncodeToString(newpk), tt.pubKey.Value)
			}
		})
	}
}

func Test_secAccnUint64(t *testing.T) {
	tests := []struct {
		name       string
		ba         *BaseAccount
		wantAccNum uint64
		wantSeq    uint64
		mustFail   bool
	}{
		{
			"empty baseaccount",
			&BaseAccount{},
			0, 0,
			true,
		},
		{
			"baseaccount only with account number",
			&BaseAccount{
				AccountNumber: "1",
			},
			0, 0,
			true,
		},
		{
			"baseaccount only with sequence number",
			&BaseAccount{
				Sequence: "1",
			},
			0, 0,
			true,
		},
		{
			"baseaccount with sequence number and account number",
			&BaseAccount{
				Sequence:      "1",
				AccountNumber: "1",
			},
			1, 1,
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.mustFail {
				require.Panics(t, func() { _, _ = secAccnUint64(tt.ba) })
			} else {
				accn, seq := secAccnUint64(tt.ba)
				require.Equal(t, tt.wantAccNum, accn)
				require.Equal(t, tt.wantSeq, seq)
			}
		})
	}
}
