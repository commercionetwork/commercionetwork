package types

import (
	"fmt"
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
		{
			"invalid serviceEndpoint",
			Service{
				ID:              "did:example:123456789abcdefghi#vcr",
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://l ocalhost",
			},
			false,
			sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"serviceEndpoint\" does not contain a valid URL"),
		},
		{
			"signatureprices defined but type not okay",
			Service{
				ID:              "did:example:123456789abcdefghi#vcr",
				Type:            "CredentialRepositoryService",
				ServiceEndpoint: "http://theUrl",
				SignaturePrices: []SignaturePrice{
					{
						CertificateProfile: "t",
						Price:              nil,
					},
				},
			},
			false,
			sdkErr.Wrap(sdkErr.ErrInvalidRequest, "signature_prices present but service type not \"signature\""),
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
			[]SignaturePrice{},
		},
		{
			"service2",
			"type2",
			"entrypoint2",
			[]SignaturePrice{},
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
					[]SignaturePrice{},
				},
				{
					"service1",
					"type1",
					"entrypoint1",
					[]SignaturePrice{},
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
					[]SignaturePrice{},
				},
				{
					"service2",
					"type2",
					"entrypoint2",
					[]SignaturePrice{},
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
				{"otherId", "otherType", "otherEndpoint", []SignaturePrice{}},
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
		{
			"document has more than 2 services with the same id",
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
				Service: Services{
					Service{
						ID:              "1",
						Type:            "1",
						ServiceEndpoint: "1",
					},
					Service{
						ID:              "1",
						Type:            "1",
						ServiceEndpoint: "1",
					},
				},
			},
			true,
		},
		{
			"document has 2 services with the same id",
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
				Service: Services{
					Service{
						ID:              "1",
						Type:            "1",
						ServiceEndpoint: "1",
					},
					Service{
						ID:              "1",
						Type:            "1",
						ServiceEndpoint: "1",
					},
					Service{
						ID:              "1",
						Type:            "1",
						ServiceEndpoint: "1",
					},
				},
			},
			true,
		},
		{
			"document has 2 services with signature types",
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
				Service: Services{
					Service{
						ID:              "3",
						Type:            "signature",
						ServiceEndpoint: "localhost",
					},
					Service{
						ID:              "1",
						Type:            "signature",
						ServiceEndpoint: "localhost",
					},
				},
			},
			true,
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

func TestSignaturePrice_Equal(t *testing.T) {
	a := sdk.NewCoin("ucommercio", sdk.NewInt(42))
	b := sdk.NewCoin("ucommercio", sdk.NewInt(43))
	tests := []struct {
		name   string
		first  SignaturePrice
		second SignaturePrice
		want   bool
	}{
		{
			"equal",
			SignaturePrice{
				CertificateProfile: "c",
				Price:              &a,
			},
			SignaturePrice{
				CertificateProfile: "c",
				Price:              &a,
			},
			true,
		},
		{
			"different cert profile",
			SignaturePrice{
				CertificateProfile: "c",
				Price:              &a,
			},
			SignaturePrice{
				CertificateProfile: "cc",
				Price:              &a,
			},
			false,
		},
		{
			"different price",
			SignaturePrice{
				CertificateProfile: "c",
				Price:              &a,
			},
			SignaturePrice{
				CertificateProfile: "c",
				Price:              &b,
			},
			false,
		},
		{
			"everything different",
			SignaturePrice{
				CertificateProfile: "c",
				Price:              &a,
			},
			SignaturePrice{
				CertificateProfile: "cc",
				Price:              &b,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.first.Equal(tt.second))
		})
	}
}

func TestSignaturePrices_Equal(t *testing.T) {
	a := sdk.NewCoin("ucommercio", sdk.NewInt(42))

	tests := []struct {
		name   string
		first  SignaturePrices
		second SignaturePrices
		want   bool
	}{
		{
			"equal",
			SignaturePrices{
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
					MembershipMultiplier: map[string]sdk.Dec{
						"price": sdk.NewDec(42),
					},
				},
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
				},
			},
			SignaturePrices{
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
					MembershipMultiplier: map[string]sdk.Dec{
						"price": sdk.NewDec(42),
					},
				},
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
				},
			},
			true,
		},
		{
			"not equal",
			SignaturePrices{
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
					MembershipMultiplier: map[string]sdk.Dec{
						"price": sdk.NewDec(42),
					},
				},
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
				},
			},
			SignaturePrices{
				SignaturePrice{
					CertificateProfile: "cc",
					Price:              &a,
				},
				SignaturePrice{
					CertificateProfile: "c",
					Price:              &a,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.first.Equal(tt.second))
		})
	}
}

func TestSignaturePrices_Price(t *testing.T) {
	p := sdk.NewCoin("ucommercio", sdk.NewInt(42))
	a := &p

	halfOff, err := sdk.NewDecFromStr("0.5")
	require.NoError(t, err)

	tests := []struct {
		name          string
		sp            SignaturePrices
		certp         string
		membership    string
		expectedPrice sdk.Coin
		wantErr       bool
	}{
		{
			"no prices defined",
			SignaturePrices{},
			"profile",
			"membership",
			sdk.NewCoin("ucommercio", sdk.NewInt(0)),
			false,
		},
		{
			"price for membership available with a 2x multiplier, but user has none",
			SignaturePrices{
				{
					CertificateProfile: "profile",
					Price:              a,
					MembershipMultiplier: map[string]sdk.Dec{
						"membership": sdk.NewDec(2),
					},
				},
			},
			"profile",
			"",
			sdk.NewCoin("ucommercio", sdk.NewInt(42)),
			false,
		},
		{
			"price for membership available with a 2x multiplier",
			SignaturePrices{
				{
					CertificateProfile: "profile",
					Price:              a,
					MembershipMultiplier: map[string]sdk.Dec{
						"membership": sdk.NewDec(2),
					},
				},
			},
			"profile",
			"membership",
			sdk.NewCoin("ucommercio", sdk.NewInt(84)),
			false,
		},
		{
			"price for membership available with a 50% multiplier",
			SignaturePrices{
				{
					CertificateProfile: "profile",
					Price:              a,
					MembershipMultiplier: map[string]sdk.Dec{
						"membership": halfOff,
					},
				},
			},
			"profile",
			"membership",
			sdk.NewCoin("ucommercio", sdk.NewInt(21)),
			false,
		},
		{
			"price available",
			SignaturePrices{
				{
					CertificateProfile: "profile",
					Price:              a,
				},
			},
			"profile",
			"membership",
			p,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, e := tt.sp.Price(tt.certp, tt.membership)
			if tt.wantErr {
				require.Nil(t, p)
				require.Error(t, e)
				return
			}

			require.NoError(t, e)

			if tt.expectedPrice.Amount.IsZero() {
				require.Nil(t, p)
			} else {
				require.True(t, tt.expectedPrice.IsEqual(*p))
			}
		})
	}
}

func TestServices_SignatureEnabled(t *testing.T) {
	tests := []struct {
		name        string
		sp          Services
		signEnabled bool
	}{
		{
			"signature enabled",
			Services{
				Service{
					ID:              "id",
					Type:            SignatureService,
					ServiceEndpoint: "localhost",
				},
			},
			true,
		},
		{
			"signature not enabled",
			Services{
				Service{
					ID:              "id",
					Type:            "nosig",
					ServiceEndpoint: "localhost",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, e := tt.sp.SignatureEnabled()
			if tt.signEnabled {
				require.NotEqual(t, Service{}, s)
				require.NoError(t, e)
				return
			}

			require.Equal(t, Service{}, s)
			require.Error(t, e)
		})
	}
}
