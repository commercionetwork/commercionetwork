package types

// Representation of a priced token of any type
type Asset struct {
	Name string
	Code string
}

func (asset Asset) Equals(a Asset) bool {
	return asset.Name == a.Name &&
		asset.Code == a.Code
}

type Assets []Asset

func (assets Assets) AppendIfMissing(a Asset) Assets {
	for _, ele := range assets {
		if ele.Equals(a) {
			return assets
		}
	}
	return append(assets, a)
}
