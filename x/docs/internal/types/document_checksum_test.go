package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentChecksum_Validate(t *testing.T) {
	tests := []struct {
		name     string
		checksum DocumentChecksum
		wantErr  string
	}{
		{
			"empty value",
			DocumentChecksum{
				Value:     "",
				Algorithm: "md5",
			},
			"checksum value can't be empty",
		},
		{
			"empty algorithm",
			DocumentChecksum{
				Value:     "48656c6c6f20476f7068657221234567",
				Algorithm: "",
			},
			"checksum algorithm can't be empty",
		},
		{
			"invalid hex value",
			DocumentChecksum{
				Value:     "qr54g7srg5674fsg4sfg",
				Algorithm: "md5",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"invalid algorithm type",
			DocumentChecksum{
				Value:     "48656c6c6f20476f7068657221234567",
				Algorithm: "md6",
			},
			"invalid checksum algorithm type md6",
		},
		{
			"bad md5 hex hash decoding",
			DocumentChecksum{
				Value:     "0cc175bc0f1b6a831c399e269772661",
				Algorithm: "md5",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"bad sha-1 hex hash decoding",
			DocumentChecksum{
				Value:     "86f7e437faa5a7fce15dddcb9eaeaea377667b8",
				Algorithm: "sha-1",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"bad sha-224 hex hash decoding",
			DocumentChecksum{
				Value:     "abd37534c7d9a2efb946de931cd7055ffdb8879563ae98078d6d6d5",
				Algorithm: "sha-224",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"bad sha-256 hex hash decoding",
			DocumentChecksum{
				Value:     "ca978112ca1bbdcafac21b39a23dc4da786eff8147c4e72b9807785afee48bb",
				Algorithm: "sha-256",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"bad sha-384 hex hash decoding",
			DocumentChecksum{
				Value:     "54a59b9f22b0b80880d427e548b7c23abd873486e1f035dce9cd697e85175033caa88e6d57bc35efae0b5afd3145f31",
				Algorithm: "sha-384",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"bad sha-512 hex hash decoding",
			DocumentChecksum{
				Value:     "1f40fc92da24169475099ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
				Algorithm: "sha-512",
			},
			"invalid checksum value (must be hex)",
		},
		{
			"good md5 hex hash decoding",
			DocumentChecksum{
				Value:     "0cc175b9c0f1b6a831c399e269772661",
				Algorithm: "md5",
			},
			"",
		},
		{
			"good sha-1 hex hash decoding",
			DocumentChecksum{
				Value:     "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8",
				Algorithm: "sha-1",
			},
			"",
		},
		{
			"good sha-224 hex hash decoding",
			DocumentChecksum{
				Value:     "abd37534c7d9a2efb9465de931cd7055ffdb8879563ae98078d6d6d5",
				Algorithm: "sha-224",
			},
			"",
		},
		{
			"good sha-256 hex hash decoding",
			DocumentChecksum{
				Value:     "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
				Algorithm: "sha-256",
			},
			"",
		},
		{
			"good sha-384 hex hash decoding",
			DocumentChecksum{
				Value:     "54a59b9f22b0b80880d8427e548b7c23abd873486e1f035dce9cd697e85175033caa88e6d57bc35efae0b5afd3145f31",
				Algorithm: "sha-384",
			},
			"",
		},
		{
			"good sha-512 hex hash decoding",
			DocumentChecksum{
				Value:     "1f40fc92da241694750979ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
				Algorithm: "sha-512",
			},
			"",
		},
		{
			"a well-defined hex value for a well-known algorithm, but with bad hash length",
			DocumentChecksum{
				Value:     "68656c6c6f2c20746573747321",
				Algorithm: "md5",
			},
			"invalid checksum length for algorithm md5",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.checksum.Validate()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDocumentChecksum_Equals(t *testing.T) {
	tests := []struct {
		name  string
		us    DocumentChecksum
		other DocumentChecksum
		equal bool
	}{
		{
			"two equal DocumentChecksums",
			DocumentChecksum{
				Value:     "0cc175b9c0f1b6a831c399e269772661",
				Algorithm: "md5",
			},
			DocumentChecksum{
				Value:     "0cc175b9c0f1b6a831c399e269772661",
				Algorithm: "md5",
			},
			true,
		},
		{
			"two different DocumentChecksums",
			DocumentChecksum{
				Value:     "0cc175b9c0f1b6a831c399e269772661",
				Algorithm: "md5",
			},
			DocumentChecksum{
				Value:     "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
				Algorithm: "sha-256",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val := tt.us.Equals(tt.other)
			assert.Equal(t, tt.equal, val)
		})
	}
}
