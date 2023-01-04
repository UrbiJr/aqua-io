package base

import "github.com/UrbiJr/anticaptcha"

// CaptchaSolver interface defines methods used by each configured ServiceProvider.
type CaptchaSolver interface {
	SolveGeetest(captcha *GeetestCaptcha) (anticaptcha.GeeTestSolution, error)
	SolveHCaptcha(captcha *HCaptcha) (string, error)
	SolveReCaptchaV2(captcha *ReCaptchaV2) (string, error)
	SolveReCaptchaV3(captcha *ReCaptchaV3) (string, error)
	SolveReCaptchaV2E(captcha *ReCaptchaV2E) (string, error)
	SolveReCaptchaV3E(captcha *ReCaptchaV3E) (string, error)
	SolveNormalCaptcha(captcha *NormalCaptcha) (string, error)
	SolveTextCaptcha(captcha *TextCaptcha) (string, error)
}

// ServiceProvider are supported solving services used to configure a CaptchaSolver
type ServiceProvider string

const (
	// TwoCaptcha (2Captcha API) - https://2captcha.com/enterpage
	TwoCaptcha ServiceProvider = "2captcha"
	// AYCD AutoSolve - https://aycd.zendesk.com/hc/en-us/articles/360043730513-What-is-AutoSolve-
	AYCD        ServiceProvider = "aycd"
	Anticaptcha ServiceProvider = "anticaptcha"
	Capmonster  ServiceProvider = "capmonster"
)

// SolverOptions defines options used to configure a CaptchaSolver
type SolverOptions struct {
	Provider    ServiceProvider
	APIKey      string
	AccessToken string
}
