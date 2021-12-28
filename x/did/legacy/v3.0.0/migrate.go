package v3_0_0

import (
	"strings"

	v220did "github.com/commercionetwork/commercionetwork/x/did/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220did.GenesisState) *types.GenesisState {

	didDocuments := DidDocuments{}
	var didDocument *types.DidDocument
	for _, v220didDocument := range oldGenState.DidDocuments {
		didDocument = convertDDO(v220didDocument)
		didDocuments = didDocuments.AppendIfMissingID(didDocument)
	}

	return &types.GenesisState{DidDocuments: didDocuments}
}

func publicKeyPemToMultiBase(pkPem string) (pkMultiBase string) {
	pkMultiBase = pkPem
	pkMultiBase = strings.ReplaceAll(pkMultiBase, "\n", "")
	pkMultiBase = strings.ReplaceAll(pkMultiBase, "\r", "")
	pkMultiBase = strings.ReplaceAll(pkMultiBase, "-", "")
	pkMultiBase = strings.TrimPrefix(pkMultiBase, "BEGIN PUBLIC KEY")
	pkMultiBase = strings.TrimSuffix(pkMultiBase, "END PUBLIC KEY")
	// add multibase code for base64 (rfc4648 no padding)
	pkMultiBase = "m" + pkMultiBase
	return
}

func convertPubKeys(pubKeys v220did.PubKeys) (verificationMethods []*types.VerificationMethod) {

	for _, pubKey := range pubKeys {

		verificationMethod := types.VerificationMethod{
			ID:                 pubKey.ID,
			Type:               pubKey.Type,
			Controller:         pubKey.Controller.String(),
			PublicKeyMultibase: publicKeyPemToMultiBase(pubKey.PublicKeyPem),
		}

		verificationMethods = append(verificationMethods, &verificationMethod)
	}

	return
}

func convertService(services220 v220did.Services) (services300 []*types.Service) {

	for _, service220 := range services220 {
		service300 := types.Service{
			ID:              service220.ID,
			Type:            service220.Type,
			ServiceEndpoint: service220.ServiceEndpoint,
		}
		services300 = append(services300, &service300)
	}

	return
}

func convertDDO(ddo220 v220did.DidDocument) (ddo300 *types.DidDocument) {

	ddo300.Context = []string{ddo220.Context}

	ddo300.ID = ddo220.ID.String()

	ddo300.VerificationMethod = convertPubKeys(ddo220.PubKeys)

	ddo300.Service = convertService(ddo220.Service)

	// created ?

	return
}
