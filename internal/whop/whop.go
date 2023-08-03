package whop

type Whop struct {
	APIBaseEndpoint string
	AuthAPIKey      string
}

// InitWhop returns a new instance of Whop helper
func InitWhop() *Whop {
	settings := &Whop{
		APIBaseEndpoint: "https://api.whop.com/api/v2/",
		AuthAPIKey:      "-0ZDuJc64nKt-PkiXa7WNkg3z_7eFh97odFor1An4kk",
	}

	return settings
}
