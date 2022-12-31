package base

const (
	HCAPTCHA      = 1
	RECAPTCHA     = 2
	RECAPTCHAV2   = 3
	RECAPTCHAV2E  = 4
	RECAPTCHAV3   = 5
	RECAPTCHAV3E  = 6
	NORMALCAPTCHA = 7
	TEXTCAPTCHA   = 8
)

// GeetestCaptcha defines parameters used to solve Geetest.
type GeetestCaptcha struct {
	PageURL      string
	ApiServer    string
	GtKey        string
	ChallengeKey string
	Proxy        string // use format "LOGIN:PASSWORD@IP:PORT"
}

// HCaptcha defines parameters used to solve HCaptcha
type HCaptcha struct {
	PageURL   string
	SiteKey   string
	Invisible bool
	Proxy     string // use format "LOGIN:PASSWORD@IP:PORT"
}

// ReCaptcha defines parameters used to solve ReCaptcha
type ReCaptcha struct {
	PageURL string
	SiteKey string
}

// ReCaptchaV2 defines parameters used to solve ReCaptcha V2
type ReCaptchaV2 struct {
	PageURL   string
	SiteKey   string
	Invisible bool
	Action    string // e.g. "verify"
	Proxy     string // use format "LOGIN:PASSWORD@IP:PORT"
}

// ReCaptchaV2E defines parameters used to solve Enterprise ReCaptcha V2
type ReCaptchaV2E struct {
	PageURL   string
	SiteKey   string
	Invisible bool
	DataS     string
}

// ReCaptchaV3 defines parameters used to solve ReCaptcha V3
type ReCaptchaV3 struct {
	PageURL  string
	SiteKey  string
	MinScore float32
	Action   string // e.g. "verify"
	Proxy    string // use format "LOGIN:PASSWORD@IP:PORT"
}

// ReCaptchaV3E defines parameters used to solve Enterprise ReCaptcha V3
type ReCaptchaV3E struct {
	PageURL  string
	SiteKey  string
	Action   string
	MinScore float32
	DataS    string
}

// NormalCaptcha defines parameters used to solve Normal Captchas
type NormalCaptcha struct {
	Language        string
	Body            []byte
	Instructions    interface{}
	Phrase          bool
	CaseSensitive   bool
	Numeric         bool
	Letters         bool
	AlphaNumericXOR bool
	AlphaNumericAND bool
	Math            bool
}

// TextCaptcha defines parameters used to solve Text Captchas
type TextCaptcha struct {
	Language string
	Text     string
}
