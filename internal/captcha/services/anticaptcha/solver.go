package anticaptcha

import (
	"errors"
	"time"

	"github.com/UrbiJr/anticaptcha"
	"github.com/UrbiJr/aqua-io/internal/captcha/base"
)

var manager = AntiCaptchaManager{}

type CaptchaSolver struct{}

func NewCaptchaSolver(options base.SolverOptions) *CaptchaSolver {
	manager.Credential.APIKey = options.APIKey
	return &CaptchaSolver{}
}

// SolveGeetest is used to get a token for either the checkbox or invisible reCaptcha V2 CAPTCHA
func (solver *CaptchaSolver) SolveGeetest(captcha *base.GeetestCaptcha) (anticaptcha.GeeTestSolution, error) {
	return manager.GetGeetest(captcha.GtKey, captcha.ChallengeKey, captcha.ApiServer, captcha.PageURL)
}

// SolveReCaptchaV2 is used to get a token for either the checkbox or invisible reCaptcha V2 CAPTCHA
func (solver *CaptchaSolver) SolveReCaptchaV2(captcha *base.ReCaptchaV2) (string, error) {
	return manager.GetRecaptchaV2(captcha.PageURL, captcha.SiteKey)
}

// SolveReCaptchaV2E is unsupported by AutoSolve but must be implemented to satisfy the base
// CaptchaSolver interface.
func (solver *CaptchaSolver) SolveReCaptchaV2E(captcha *base.ReCaptchaV2E) (string, error) {
	return "", ErrCaptchaUnsupported
}

// SolveReCaptchaV3 is used to get a token for the normal reCaptcha V3 CAPTCHA
func (solver *CaptchaSolver) SolveReCaptchaV3(captcha *base.ReCaptchaV3) (string, error) {
	return "", ErrCaptchaUnsupported
}

// SolveReCaptchaV3E is used to get a token for the enterprise reCaptcha V3 CAPTCHA
func (solver *CaptchaSolver) SolveReCaptchaV3E(captcha *base.ReCaptchaV3E) (string, error) {
	return "", ErrCaptchaUnsupported
}

// SolveHCaptcha is used to get a token for either the checkbox or invisible hCaptcha CAPTCHA
func (solver *CaptchaSolver) SolveHCaptcha(captcha *base.HCaptcha) (string, error) {
	taskId, err := manager.GetHCaptchaCreateTask(captcha.PageURL, captcha.SiteKey)

	if err != nil {
		return "", err
	}

	check := time.NewTicker(10 * time.Second)
	timeout := time.NewTimer(150 * time.Second)

	for {
		select {
		case <-check.C:
			response, err := manager.getTaskResult(taskId)
			if err != nil {
				return "", err
			}
			if response["status"] == "ready" {
				return response["solution"].(map[string]any)["gRecaptchaResponse"].(string), nil
			}
			check = time.NewTicker(checkInterval)
		case <-timeout.C:
			return "", errors.New("antiCaptcha check result timeout")
		}
	}
}

// SolveNormalCaptcha is unsupported by AutoSolve but must be implemented to satisfy the base
// CaptchaSolver interface.
func (solver *CaptchaSolver) SolveNormalCaptcha(captcha *base.NormalCaptcha) (string, error) {
	return "", ErrCaptchaUnsupported
}

// SolveTextCaptcha is unsupported by AutoSolve but must be implemented to satisfy the base
// CaptchaSolver interface.
func (solver *CaptchaSolver) SolveTextCaptcha(captcha *base.TextCaptcha) (string, error) {
	return "", ErrCaptchaUnsupported
}
