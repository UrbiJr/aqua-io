package sites

type Store string
type SiteCategory int
type MenuPage int

const (

	// Main Modules:

	Nike  Store = "Nike"
	Schuh Store = "Schuh"
)

func (s Store) String() string {
	return string(s)
}

const (
	SneakerSite SiteCategory = iota
	NonSneakerSite
)

// SupportedSite contains information about a supported module
type SupportedSite struct {
	Name      Store
	Category  SiteCategory
	CSVFields []string
}
