package types_test

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestProof_Equals(t *testing.T) {
	zone, _ := time.LoadLocation("UTC")
	date := time.Date(2019, 1, 1, 1, 1, 1, 1, zone)
	proof := types.NewProof("type-1", date, "creator-1", "")

	assert.False(t, proof.Equals(types.NewProof("type-2", proof.Created, proof.Creator, proof.SignatureValue)))
	assert.False(t, proof.Equals(types.NewProof(proof.Type, proof.Created.AddDate(0, 0, 1), proof.Creator, proof.SignatureValue)))
	assert.False(t, proof.Equals(types.NewProof(proof.Type, proof.Created, "creator-2", proof.SignatureValue)))
	assert.False(t, proof.Equals(types.NewProof(proof.Type, proof.Created, proof.Creator, proof.SignatureValue+"1")))
	assert.True(t, proof.Equals(proof))
}

func TestProof_Validate(t *testing.T) {
	zone, _ := time.LoadLocation("UTC")
	date := time.Date(2019, 1, 1, 1, 1, 1, 1, zone)

	assert.Error(t, types.NewProof("", date, "creator", "signature").Validate())
	assert.Error(t, types.NewProof("LinkedDataSignature2015", time.Time{}, "creator", "signature").Validate())
	assert.NoError(t, types.NewProof("LinkedDataSignature2015", date, "creator", "signatureValue").Validate())
}
