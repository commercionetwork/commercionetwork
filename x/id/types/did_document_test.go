package types

import (
	"fmt"
	"strings"
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
		Proof: &Proof{
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

func TestDidDocument_VerifyProof(t *testing.T) {

	var testZone, _ = time.LoadLocation("UTC")
	var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)
	var testOwnerAddress, _ = sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	okayDDOWrongSig := DidDocument{
		Context: ContextDidV1,
		ID:      testOwnerAddress,
		Proof: &Proof{
			Type:               KeyTypeSecp256k12019,
			Created:            testTime,
			ProofPurpose:       ProofPurposeAuthentication,
			Controller:         testOwnerAddress.String(),
			SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
			VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
		},
		PubKeys: PubKeys{
			PubKey{
				ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
				Type:       "RsaVerificationKey2018",
				Controller: testOwnerAddress,
				PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
			},
			PubKey{
				ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
				Type:       "RsaSignatureKey2018",
				Controller: testOwnerAddress,
				PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
			},
		},
	}

	okayDDO := DidDocument{
		Context: ContextDidV1,
		ID:      testOwnerAddress,
		Proof: &Proof{
			Type:               KeyTypeSecp256k12019,
			Created:            testTime,
			ProofPurpose:       ProofPurposeAuthentication,
			Controller:         testOwnerAddress.String(),
			SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
			VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
		},
		PubKeys: PubKeys{
			PubKey{
				ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
				Type:       "RsaVerificationKey2018",
				Controller: testOwnerAddress,
				PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
			},
			PubKey{
				ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
				Type:       "RsaSignatureKey2018",
				Controller: testOwnerAddress,
				PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
			},
		},
	}

	tests := []struct {
		name        string
		didDocument DidDocument
		wantErr     bool
	}{
		{
			"pubkey is not bech32",
			DidDocument{
				Proof: &Proof{
					VerificationMethod: "notbech32",
				},
			},
			true,
		},
		{
			"signaturevalue is not base64",
			DidDocument{
				Proof: &Proof{
					VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
					SignatureValue:     "notbase64",
				},
			},
			true,
		},
		{
			"signature wrong",
			okayDDOWrongSig,
			true,
		},
		{
			"signature ok, no errors",
			okayDDO,
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.didDocument.VerifyProof()

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDidDocument_Validate(t *testing.T) {
	var testOwnerAddress, _ = sdk.AccAddressFromBech32("did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf")
	var testZone, _ = time.LoadLocation("UTC")
	var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)

	tests := []struct {
		name        string
		didDocument DidDocument
		wantErr     bool
	}{
		{
			"empty id",
			DidDocument{},
			true,
		},
		{
			"context is not ContextDidV1",
			DidDocument{
				Context: "ohno",
			},
			true,
		},
		{
			"no pubkeys",
			DidDocument{
				ID:      testOwnerAddress,
				Context: ContextDidV1,
			},
			true,
		},
		{
			"no required pubkeys",
			DidDocument{
				ID:      testOwnerAddress,
				Context: ContextDidV1,
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
			},
			true,
		},
		{
			"no required pubkeys",
			DidDocument{
				ID:      testOwnerAddress,
				Context: ContextDidV1,
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
			},
			true,
		},
		{
			"no proof",
			DidDocument{
				ID:      testOwnerAddress,
				Context: ContextDidV1,
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
			},
			true,
		},
		{
			"service invalid",
			DidDocument{
				ID:      testOwnerAddress,
				Context: ContextDidV1,
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
				Proof: &Proof{
					Type:               KeyTypeSecp256k12019,
					Created:            testTime,
					ProofPurpose:       ProofPurposeAuthentication,
					Controller:         testOwnerAddress.String(),
					SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
					VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
				},
				Service: Services{
					Service{},
				},
			},
			true,
		},
		{
			"duplicate keys (2 keys)",
			DidDocument{
				Context: ContextDidV1,
				ID:      testOwnerAddress,
				Proof: &Proof{
					Type:               KeyTypeSecp256k12019,
					Created:            testTime,
					ProofPurpose:       ProofPurposeAuthentication,
					Controller:         testOwnerAddress.String(),
					SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
					VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
				},
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
			},
			true,
		},
		{
			"duplicate keys (more than 2 keys)",
			DidDocument{
				Context: ContextDidV1,
				ID:      testOwnerAddress,
				Proof: &Proof{
					Type:               KeyTypeSecp256k12019,
					Created:            testTime,
					ProofPurpose:       ProofPurposeAuthentication,
					Controller:         testOwnerAddress.String(),
					SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
					VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
				},
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-3", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
			},
			true,
		},
		{
			"all ok",
			DidDocument{
				Context: ContextDidV1,
				ID:      testOwnerAddress,
				Proof: &Proof{
					Type:               KeyTypeSecp256k12019,
					Created:            testTime,
					ProofPurpose:       ProofPurposeAuthentication,
					Controller:         testOwnerAddress.String(),
					SignatureValue:     "uv9ZM4XusZl2q6Ei2O7aZW32pzwfg6ZQpBsQPb8cxzlFXWEyZLxem29fQBB4Py3W5gaXFEyPGruMXNsNDnr4sQ==",
					VerificationMethod: "did:com:pub1addwnpepqwzc44ggn40xpwkfhcje9y7wdz6sunuv2uydxmqjrvcwff6npp2exy5dn6c",
				},
				PubKeys: PubKeys{
					PubKey{
						ID:         fmt.Sprintf("%s#keys-1", testOwnerAddress),
						Type:       "RsaVerificationKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch
2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1Bh
Co06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIM
V1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPic
bLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1W
gNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3aw
GwIDAQAB
-----END PUBLIC KEY-----`,
					},
					PubKey{
						ID:         fmt.Sprintf("%s#keys-2", testOwnerAddress),
						Type:       "RsaSignatureKey2018",
						Controller: testOwnerAddress,
						PublicKeyPem: `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+Juw6xqYchTNFYUznmoB
CzKfQG75v2Pv1Db1Z5EJgP6i0yRsBG1VqIOY4icRnyhDDVFi1omQjjUuCRxWGjsc
B1UkSnybm0WC+g82HL3mUzbZja27NFJPuNaMaUlNbe0daOG88FS67jq5J2LsZH/V
cGZBX5bbtCe0Niq39mQdJxdHq3D5ROMA73qeYvLkmXS6Dvs0w0fHsy+DwJtdOnOj
xt4F5hIEXGP53qz2tBjCRL6HiMP/cLSwAd7oc67abgQxfnf9qldyd3X0IABpti1L
irJNugfN6HuxHDm6dlXVReOhHRbkEcWedv82Ji5d/sDZ+WT+yWILOq03EJo/LXJ1
SQIDAQAB
-----END PUBLIC KEY-----`,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.didDocument.Validate()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDidDocument_lengthLimits(t *testing.T) {
	bigString := strings.Repeat("c", 513)

	tests := []struct {
		name    string
		doc     DidDocument
		wantErr bool
	}{
		{
			"service id longer than 64 byte",
			DidDocument{
				Service: []Service{
					{
						ID: bigString,
					},
				},
			},
			true,
		},
		{
			"service type longer than 64 byte",
			DidDocument{
				Service: []Service{
					{
						Type: bigString,
					},
				},
			},
			true,
		},
		{
			"service endpoint longer than 512 byte",
			DidDocument{
				Service: []Service{
					{
						ServiceEndpoint: bigString,
					},
				},
			},
			true,
		},
		{
			"all fine",
			DidDocument{
				Service: []Service{
					{
						Type:            strings.Repeat("c", 64),
						ServiceEndpoint: strings.Repeat("c", 512),
						ID:              strings.Repeat("c", 64),
					},
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				require.Error(t, tt.doc.lengthLimits())
				return
			}
			require.NoError(t, tt.doc.lengthLimits())
		})
	}
}
