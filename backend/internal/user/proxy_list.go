package user

type ProxyList struct {
	ID      int64    `json:"id"`
	Title   string   `json:"title"`
	Proxies []string `json:"proxies"`
}
