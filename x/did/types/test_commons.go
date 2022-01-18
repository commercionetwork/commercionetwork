package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// package initialization for correct validation of commercionetwork addresses
func init() {
	configTestPrefixes()
}

func configTestPrefixes() {
	AccountAddressPrefix := "did:com:"
	AccountPubKeyPrefix := AccountAddressPrefix + "pub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.Seal()
}

const (
	validDidSubject   = "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc"
	validDidNoSubject = "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"

	validBase64RSAKey   = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
	invalidBase64RSAKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9ME/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
)

var (
	validContext = []string{
		ContextDidV1,
		"https://w3id.org/security/suites/ed25519-2018/v1",
		"https://w3id.org/security/suites/x25519-2019/v1",
	}

	validVerificationMethodRsaVerificationKey2018 = VerificationMethod{
		ID:                 validDidSubject + RsaVerificationKey2018NameSuffix,
		Type:               RsaVerificationKey2018,
		Controller:         validDidSubject,
		PublicKeyMultibase: string(MultibaseCodeBase64) + validBase64RSAKey,
	}

	validVerificationMethodRsaSignature2018 = VerificationMethod{
		ID:                 validDidSubject + RsaSignature2018NameSuffix,
		Type:               RsaSignature2018,
		Controller:         validDidSubject,
		PublicKeyMultibase: validVerificationMethodRsaVerificationKey2018.PublicKeyMultibase,
	}

	validVerificationMethods = []*VerificationMethod{
		&validVerificationMethodRsaVerificationKey2018,
		&validVerificationMethodRsaSignature2018,
	}

	validServiceBar = Service{
		ID:              "https://bar.example.com",
		Type:            "agent",
		ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
	}

	validServiceFoo = Service{
		ID:              "https://foo.example.com",
		Type:            "xdi",
		ServiceEndpoint: "https://commerc.io/xdi/serviceEndpoint/",
	}

	validServices = []*Service{
		&validServiceBar,
		&validServiceFoo,
	}

	validDidDocument = DidDocument{
		Context:            validContext,
		ID:                 validDidSubject,
		VerificationMethod: validVerificationMethods,
		Authentication: []string{
			validDidSubject + RsaVerificationKey2018NameSuffix,
		},
		AssertionMethod: []string{
			validDidSubject + RsaSignature2018NameSuffix,
		},
		KeyAgreement: []string{
			RsaVerificationKey2018NameSuffix,
		},
		CapabilityInvocation: []string{
			RsaSignature2018NameSuffix,
		},
		CapabilityDelegation: nil,
		Service:              validServices,
	}

	validMsgSetDidDocument = MsgSetDidDocument{
		DidDocument: &validDidDocument,
	}

	validDateCreated = "2019-03-23T06:35:22Z"
	validDateUpdated = "2023-08-10T13:40:06Z"

	validMetadata = Metadata{
		Created: validDateCreated,
		Updated: validDateUpdated,
	}

	ValidIdentity = Identity{
		DidDocument: &validDidDocument,
		Metadata:    &validMetadata,
	}
)
