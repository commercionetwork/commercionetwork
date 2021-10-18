package v3_0_0

import (
	v220did "github.com/commercionetwork/commercionetwork/x/did/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220did.GenesisState) *types.GenesisState {

	//documents
	didDocuments := DidDocuments{}
	var didDocument *types.DidDocument
	for _, v220didDocument := range oldGenState.DidDocuments {
		didDocument = migrateDidDocuments(v220didDocument)
		//documents =  append(documents, document)
		didDocuments = didDocuments.AppendIfMissingID(didDocument)
	}

	return &types.GenesisState{DidDocuments: didDocuments}

}

// migrateDidDocuments migrates a single v2.2.0 document into a 3.0.0 document
func migrateDidDocuments(didDoc v220did.DidDocument) *types.DidDocument {
	// Convert the public keys
	var pubKeys types.PubKeys
	if didDoc.PubKeys != nil {
		for _, pubKey := range didDoc.PubKeys {
			pubKeyV300 := types.PubKey{
				ID:           pubKey.ID,
				Type:         pubKey.Type,
				Controller:   pubKey.Controller.String(),
				PublicKeyPem: pubKey.PublicKeyPem,
			}
			pubKeys = append(pubKeys, &pubKeyV300)

		}
	}
	// Convert the service
	var services []*types.Service
	if didDoc.Service != nil {
		for _, service := range didDoc.Service {
			serviceV300 := types.Service{
				ID:              service.ID,
				Type:            service.Type,
				ServiceEndpoint: service.ServiceEndpoint,
			}
			services = append(services, &serviceV300)

		}
	}

	// Return a new did document
	return &types.DidDocument{
		Context: didDoc.Context,
		ID:      didDoc.ID.String(),
		PubKeys: pubKeys,
		Service: services,
	}
}
