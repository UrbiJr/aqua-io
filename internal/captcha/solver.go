package captcha

import (
	"fmt"

	"github.com/UrbiJr/nyx/internal/captcha/base"
	twocaptcha "github.com/UrbiJr/nyx/internal/captcha/services/2captcha"
	"github.com/UrbiJr/nyx/internal/captcha/services/anticaptcha"
	"github.com/UrbiJr/nyx/internal/captcha/services/capmonster"
)

type (
	// Solver provides methods used to solve captchas from a supported service provider.
	Solver = base.CaptchaSolver
	// SolverOptions defines options used to configure a CaptchaSolver.
	SolverOptions = base.SolverOptions
	// ServiceProvider are supported solving services used to configure a CaptchaSolver.
	ServiceProvider = base.ServiceProvider

	// TwoCaptchaSolver (2Captcha API) - https://2captcha.com/enterpage
	TwoCaptchaSolver = twocaptcha.CaptchaSolver

	// Anticaptcha Solver
	AnticaptchaSolver = anticaptcha.CaptchaSolver

	// Capmonster Solver
	CapmonsterSolver = capmonster.CaptchaSolver

	// GeetestCaptcha defines parameters used to solve geetest.
	GeetestCaptcha = base.GeetestCaptcha
	// HCaptcha defines parameters used to solve HCaptcha.
	HCaptcha = base.HCaptcha
	// ReCaptchaV2 defines parameters used to solve ReCaptcha V2.
	ReCaptchaV2 = base.ReCaptchaV2
	// ReCaptchaV2E defines parameters used to solve Enterprise ReCaptcha V2.
	ReCaptchaV2E = base.ReCaptchaV2E
	// ReCaptchaV3 defines parameters used to solve ReCaptcha V3.
	ReCaptchaV3 = base.ReCaptchaV3
	// ReCaptchaV3E defines parameters used to solve Enterprise ReCaptcha V3.
	ReCaptchaV3E = base.ReCaptchaV3E

	// TextCaptcha defines parameters used to solve Text Captchas.
	TextCaptcha = base.TextCaptcha
	// NormalCaptcha defines parameters used to solve Normal Captchas.
	NormalCaptcha = base.NormalCaptcha
)

const (
	// TwoCaptcha (2Captcha API) - https://2captcha.com/enterpage
	TwoCaptcha = base.TwoCaptcha
	// AYCD AutoSolve - https://aycd.zendesk.com/hc/en-us/articles/360043730513-What-is-AutoSolve-
	AYCD        = base.AYCD
	Anticaptcha = base.Anticaptcha
	Capmonster  = base.Capmonster
)

// NewCaptchaSolver returns a new CaptchaSolver configured to a supported ServiceProvider.
func NewCaptchaSolver(options SolverOptions) (Solver, error) {
	switch options.Provider {
	case TwoCaptcha:
		return twocaptcha.NewCaptchaSolver(options), nil
	case Anticaptcha:
		return anticaptcha.NewCaptchaSolver(options), nil
	case Capmonster:
		return capmonster.NewCaptchaSolver(options), nil
	default:
		return nil, fmt.Errorf("Invalid captcha solving service provider: %q", options.Provider)
	}
}
