package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// package initialization for correct validation of commercionetwork addresses
func init() {
	ConfigTestPrefixes()
}

func ConfigTestPrefixes() {
	AccountAddressPrefix := "did:com:"
	AccountPubKeyPrefix := AccountAddressPrefix + "pub"
	ValidatorAddressPrefix := AccountAddressPrefix + "valoper"
	ValidatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	ConsNodeAddressPrefix := AccountAddressPrefix + "valcons"
	ConsNodePubKeyPrefix := AccountAddressPrefix + "valconspub"
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(ValidatorAddressPrefix, ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(ConsNodeAddressPrefix, ConsNodePubKeyPrefix)
	config.Seal()
}

const didSubject = "did:com:14zk9u8894eg7fhgw0dsesnqzmlrx85ga9rvnjc"
const didNoSubject = "did:com:18h03de6awcjk4u9gaz8s5l0xxl8ulxjctzsytd"

var validContext = []string{
	ContextDidV1,
	"https://w3id.org/security/suites/ed25519-2018/v1",
	"https://w3id.org/security/suites/x25519-2019/v1",
}

const validBase64RSAKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9MELZbOWd9xkKlPicbLp8coZh9NgLajMDWKXpuHQ8cdJSxQ/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"
const invalidBase64RSAKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAqOoLR843vgkFGudQsjch2K85QJ4Hh7l2jjrMesQFDWVcW1xr//eieGzxDogWx7tMOtQ0hw77NAURhldek1BhCo06790YHAE97JqgRQ+IR9Dl3GaGVQ2WcnknO4B1cvTRJmdsqrN1Bs4Qfd+jjKIMV1tz8zU9NmdR+DvGkAYYxoIx74YaTAxH+GCArfWMG1tRJPI9ME/ekZaTuEy7qbjbGBMVzbjhPjcxffQmGV1WgNY1BGplZz9mbBmH7siKnKIVZ5Bp55uLfEw+u2yOVx/0yKUdsmZoe4jhevCSq3awGwIDAQAB"

var validVerificationMethodRsaVerificationKey2018 = VerificationMethod{
	ID:                 didSubject + RsaVerificationKey2018NameSuffix,
	Type:               RsaVerificationKey2018,
	Controller:         didSubject,
	PublicKeyMultibase: string(MultibaseCodeBase64) + validBase64RSAKey,
}

var validVerificationMethodRsaSignature2018 = VerificationMethod{
	ID:                 didSubject + RsaSignature2018NameSuffix,
	Type:               RsaSignature2018,
	Controller:         didSubject,
	PublicKeyMultibase: validVerificationMethodRsaVerificationKey2018.PublicKeyMultibase,
}

var validVerificationMethods = []*VerificationMethod{
	&validVerificationMethodRsaVerificationKey2018,
	&validVerificationMethodRsaSignature2018,
}

var validServiceBar = Service{
	ID:              "https://bar.example.com",
	Type:            "agent",
	ServiceEndpoint: "https://commerc.io/agent/serviceEndpoint/",
}

var validServiceFoo = Service{
	ID:              "https://foo.example.com",
	Type:            "xdi",
	ServiceEndpoint: "https://commerc.io/xdi/serviceEndpoint/",
}

var validServices = []*Service{
	&validServiceBar,
	&validServiceFoo,
}

var ValidMsgSetDidDocument = MsgSetDidDocument{
	Context:            validContext,
	ID:                 didSubject,
	VerificationMethod: validVerificationMethods,
	Authentication: []string{
		didSubject + RsaVerificationKey2018NameSuffix,
	},
	AssertionMethod: []string{
		didSubject + RsaSignature2018NameSuffix,
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
