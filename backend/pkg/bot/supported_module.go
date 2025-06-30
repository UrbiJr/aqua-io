package bot

type Module string
type ModuleCategory int
type MenuPage int

const (
	// Main Modules:
	PopMart Module = "Pop Mart"
)

func (s Module) String() string {
	return string(s)
}

const (
	SneakerSite ModuleCategory = iota
	NonSneakerSite
)

// SupportedSite contains information about a supported module
type SupportedModule struct {
	Name      Module
	Category  ModuleCategory
	CSVFields []string
}

func AllModules() map[string]string {
	return map[string]string {
		"popmart" : PopMart.String(),
	}
}