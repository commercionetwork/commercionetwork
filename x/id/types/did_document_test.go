package types

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestService_Validate(t *testing.T) {
	tests := []struct {
		name        string
		theService  Service
		isValid     bool
		expectedErr error
	}{
		{
			"valid",
			Service{
				ID:              "did:example:123456789abcdefghi#vcr",
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
			},
			true,
			nil,
		},
		{
			"no id",
			Service{
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
			},
			false,
			sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "service field \"%s\" is required", "id"),
		},
		{
			"no type",
			Service{
				ID:              "did:example:123456789abcdefghi#vcr",
				ServiceEndpoint: "http://theUrl",
			},
			false,
			sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "service field \"%s\" is required", "type"),
		},
		{
			"no service point",
			Service{
				ID:   "did:example:123456789abcdefghi#vcr",
				Type: "CredentialRepositoryService",
			},
			false,
			sdkErr.Wrapf(sdkErr.ErrInvalidRequest, "service field \"%s\" is required", "serviceEndpoint"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.theService.Validate()
			if tt.isValid {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.expectedErr.Error())
		})
	}
}

func TestService_Equals(t *testing.T) {
	baseValue := Service{
		ID:              "did:example:123456789abcdefghi#vcr",
		Type:            "CredentialRepositoryService",
		ServiceEndpoint: "http://theUrl",
	}

	tests := []struct {
		name     string
		areEqual bool
		service  Service
	}{
		{
			"are equal",
			true,
			baseValue,
		},
		{
			"not equal #1: id",
			false,
			Service{
				ID:              "did:example:otherid#vcr",
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
			},
		},
		{
			"not equal #2: type",
			false,
			Service{
				ID:              "did:example:123456789abcdefghi#vcr",
				Type:            "OtherCredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
			},
		},
		{
			"not equal #3: service endpoint",
			false,
			Service{
				ID:              "did:example:123456789abcdefghi#vcr",
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://otherUrl",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.areEqual, baseValue.Equals(tt.service))
		})
	}
}

func TestServices_Equal(t *testing.T) {
	baseServices := Services{
		{
			"service1",
			"type1",
			"entrypoint1",
		},
		{
			"service2",
			"type2",
			"entrypoint2",
		},
	}
	tests := []struct {
		name     string
		services Services
		areEqual bool
	}{
		{
			"equals",
			baseServices,
			true,
		},
		{
			"different in number",
			Services{
				Service{
					ID:              "otherId1",
					Type:            "otherType1",
					ServiceEndpoint: "otherEndpoint1",
				},
			},
			false,
		},
		{
			"different in order",
			Services{
				{
					"service2",
					"type2",
					"entrypoint2",
				},
				{
					"service1",
					"type1",
					"entrypoint1",
				},
			},
			false,
		},
		{
			"value different",
			Services{
				{
					"service1",
					"type1",
					"otherEntrypoint", // The different
				},
				{
					"service2",
					"type2",
					"entrypoint2",
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.areEqual, tt.services.Equals(baseServices))
		})
	}
}

func TestDidDocument_Equals_Service(t *testing.T) {
	baseDidDocument := getBaseDocumentWithServices([]Service{
		{
			ID:              "serviceId",
			Type:            "theServiceType",
			ServiceEndpoint: "theServiceEndpoint",
		},
	})

	tests := []struct {
		name     string
		areEqual bool
		document DidDocument
	}{
		{
			"equal services",
			true,
			baseDidDocument,
		},
		{
			"different services",
			false,
			getBaseDocumentWithServices(Services{
				{"otherId", "otherType", "otherEndpoint"},
			}),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.areEqual, baseDidDocument.Equals(tt.document))
		})
	}
}

func getBaseDocumentWithServices(services []Service) DidDocument {
	addr, _ := sdk.AccAddressFromBech32("did:com:1lwmppctrr6ssnrmuyzu554dzf50apkfv6l2exu")

	return DidDocument{
		Context: "",
		ID:      addr,
		PubKeys: PubKeys{
			PubKey{
				ID:           "abcd",
				Type:         "thePubKeyType",
				Controller:   addr,
				PublicKeyPem: "thePublicKey",
			},
		},
		Proof: Proof{
			Type:               "someType",
			Created:            time.Time{},
			ProofPurpose:       "theProofPurpose",
			Controller:         "TheController",
			VerificationMethod: "verificationMethod",
			SignatureValue:     "signatureValue",
		},
		Service: services,
	}
}
