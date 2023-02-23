package user

type ProxyGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Profile contains information specific to a single account of a particular site i.e. BSTN
type Proxy struct {
	ID       int64  `json:"id"`
	GroupID  int64  `json:"group_id"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ProxyManager struct {
	Groups  []ProxyGroup
	Proxies []Proxy
}

func (pxm *ProxyManager) FilterByGroupName(groupName string) []Proxy {
	var filtered []Proxy

	// iterate through proxy groups
	for _, g := range pxm.Groups {
		// if group is the one requested
		if g.Name == groupName {
			// keep it as current and iterate through proxies
			for _, p := range pxm.Proxies {
				// if proxy belongs to current group
				if p.GroupID == g.ID {
					// add it to filtered
					filtered = append(filtered, p)
				}
			}
			return filtered
		}
	}

	return filtered
}
