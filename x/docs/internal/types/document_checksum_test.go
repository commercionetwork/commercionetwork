package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// --------------------------
// --- Checksum validation
// --------------------------

func TestValidateChecksum_validChecksum(t *testing.T) {
	var checksumList = map[string]string{
		"md5":     "0cc175b9c0f1b6a831c399e269772661",
		"sha-1":   "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8",
		"sha-224": "abd37534c7d9a2efb9465de931cd7055ffdb8879563ae98078d6d6d5",
		"sha-256": "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"sha-384": "54a59b9f22b0b80880d8427e548b7c23abd873486e1f035dce9cd697e85175033caa88e6d57bc35efae0b5afd3145f31",
		"sha-512": "1f40fc92da241694750979ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
	}

	for key, value := range checksumList {
		checksum := DocumentChecksum{Algorithm: key, Value: value}
		actual := checksum.Validate()
		assert.Nil(t, actual)
	}
}

func TestValidateChecksum_emptyValue(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "",
		Algorithm: "md5",
	}

	actual := invalidChecksum.Validate()
	assert.NotNil(t, actual)
}

func TestValidateChecksum_emptyAlgorithm(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "",
	}

	actual := invalidChecksum.Validate()
	assert.NotNil(t, actual)
}

func TestValidateChecksum_invalidHexValue(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "qr54g7srg5674fsg4sfg",
		Algorithm: "md5",
	}

	actual := invalidChecksum.Validate()
	assert.NotNil(t, actual)
}

func TestValidateChecksum_invalidAlgorithmType(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md6",
	}

	actual := invalidChecksum.Validate()
	assert.NotNil(t, actual)
}

func TestValidateChecksum_invalidChecksumLengths(t *testing.T) {
	var checksumList = map[string]string{
		"md5":     "0cc175bc0f1b6a831c399e269772661",
		"sha-1":   "86f7e437faa5a7fce15dddcb9eaeaea377667b8",
		"sha-224": "abd37534c7d9a2efb946de931cd7055ffdb8879563ae98078d6d6d5",
		"sha-256": "ca978112ca1bbdcafac21b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"sha-384": "54a59b9f22b0b80880d427e548b7c23abd873486e1f035dce9cd697e85175033caa88e6d57bc35efae0b5afd3145f31",
		"sha-512": "1f40fc92da24169475099ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
	}

	for key, value := range checksumList {
		checksum := DocumentChecksum{Algorithm: key, Value: value}
		actual := checksum.Validate()
		assert.NotNil(t, actual)
	}
}
