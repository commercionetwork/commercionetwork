package v3_0_0

import (
	"github.com/commercionetwork/commercionetwork/x/did/types"
)

const (
	ModuleName = "did"
)

type DidDocuments []*types.DidDocument

func (didDocuments DidDocuments) appendIfMissingID(i *types.DidDocument) DidDocuments {
	for _, ele := range didDocuments {
		if ele.ID == i.ID {
			return didDocuments
		}
	}
	return append(didDocuments, i)
}
