package v3_0_0

import (
	"github.com/commercionetwork/commercionetwork/x/documents/types"
)

const (
	ModuleName = types.ModuleName
)

type Documents []*types.Document

func (documents Documents) AppendIfMissingID(i *types.Document) Documents {
	for _, ele := range documents {
		if ele.UUID == i.UUID {
			return documents
		}
	}
	return append(documents, i)
}
