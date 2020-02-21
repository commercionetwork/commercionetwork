package v2_0_0

import (
	"encoding/base64"
	"fmt"
	"strconv"

	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// getPk converts our Pubkey old serialization format into a proper Tendermint crypto.Pubkey.
func getPk(p PubKey) crypto.PubKey {
	var pk crypto.PubKey

	if p.Value != "" {
		// decode public key base64
		pkHex, err := base64.StdEncoding.DecodeString(p.Value)
		if err != nil {
			panic(err)
		}

		// copy pkHex into castPk for type coherence
		castPk := secp256k1.PubKeySecp256k1{}
		copy(castPk[:], pkHex)

		pk, err = tmamino.PubKeyFromBytes(castPk.Bytes())
		if err != nil {
			panic(err)
		}
	}

	return pk
}

// secAccnUint64 converts sequence number and account number strings to their uint64 values.
func secAccnUint64(b *BaseAccount) (accNum uint64, seq uint64) {
	var err error

	accNum, err = strconv.ParseUint(b.AccountNumber, 10, 64)
	if err != nil {
		panic(fmt.Errorf("could not convert account_number to uint64: %w", err))
	}

	seq, err = strconv.ParseUint(b.Sequence, 10, 64)
	if err != nil {
		panic(fmt.Errorf("could not convert sequence to uint64: %w", err))
	}

	return
}
