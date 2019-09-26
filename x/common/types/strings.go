package types

func (elements Strings) AppendIfMissing(element string) (Strings, bool) {
	for _, ele := range elements {
		if ele == element {
			return nil, true
		}
	}
	return append(elements, element), false

}

type Strings []string
