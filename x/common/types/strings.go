package types

type Strings []string

func (elements Strings) AppendIfMissing(element string) (Strings, bool) {
	for _, ele := range elements {
		if ele == element {
			return elements, false
		}
	}
	return append(elements, element), true

}
