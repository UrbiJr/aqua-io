package twocaptcha

import (
	"encoding/json"

	api2captcha "github.com/2captcha/2captcha-go"
	"github.com/UrbiJr/anticaptcha"
	"github.com/UrbiJr/go-cactus/internal/captcha/base"
	"github.com/UrbiJr/go-cactus/internal/captcha/services/2captcha/task"
)

// CaptchaSolver client used to perform 2Captcha tasks.
type CaptchaSolver struct {
	DefaultParams *task.DefaultParams
	Client        *api2captcha.Client
}

// NewCaptchaSolver returns a new instance of CaptchaSolver
func NewCaptchaSolver(options base.SolverOptions) *CaptchaSolver {
	client := api2captcha.NewClient(options.APIKey)
	// // identifier of this software (obtained after publishing to https://2captcha.com/software/add)
	//client.SoftId = 7834

	// Timeout in seconds for all captcha types except ReCaptcha. Defines how long the module tries to get the answer from res.php API endpoint
	client.DefaultTimeout = 120

	// Timeout for ReCaptcha in seconds. Defines how long the module tries to get the answer from res.php API endpoint
	client.RecaptchaTimeout = 600

	// Interval in seconds between requests to res.php API endpoint, setting values less than 5 seconds is not recommended
	client.PollingInterval = 10

	return &CaptchaSolver{
		Client: client,
	}
}

func (solver *CaptchaSolver) SubmitRequest(request api2captcha.Request) (string, error) {

	code, err := solver.Client.Solve(request)

	if err != nil {
		if err == api2captcha.ErrTimeout {
			return "", err
		} else if err == api2captcha.ErrApi {
			return "", err
		} else if err == api2captcha.ErrNetwork {
			return "", err
		} else {
			return "", err
		}
	}
	return code, nil

}

// SolveNormalCaptcha requests a Normal Captcha solution.
func (solver *CaptchaSolver) SolveNormalCaptcha(captcha *base.NormalCaptcha) (string, error) {

	normalCaptchaTaskParams := task.NewNormalCaptchaTaskParams(captcha, solver.DefaultParams)

	cap := api2captcha.Normal{
		Base64:          normalCaptchaTaskParams.Body,
		Numberic:        normalCaptchaTaskParams.Numeric,
		Phrase:          normalCaptchaTaskParams.Phrase,
		CaseSensitive:   normalCaptchaTaskParams.CaseSensitive,
		Lang:            normalCaptchaTaskParams.Language,
		HintImageBase64: normalCaptchaTaskParams.InstructionsBody,
		HintText:        normalCaptchaTaskParams.InstructionsText,
		Calc:            normalCaptchaTaskParams.Math,
	}
	return solver.SubmitRequest(cap.ToRequest())

}

// SolveTextCaptcha requests a Text Captcha solution.
func (solver *CaptchaSolver) SolveTextCaptcha(captcha *base.TextCaptcha) (string, error) {

	cap := api2captcha.Text{
		Text: captcha.Text,
		Lang: captcha.Language,
	}
	return solver.SubmitRequest(cap.ToRequest())
}

// SolveHCaptcha requests a HCaptcha solution.
func (solver *CaptchaSolver) SolveHCaptcha(captcha *base.HCaptcha) (string, error) {
	cap := api2captcha.HCaptcha{
		SiteKey: captcha.SiteKey,
		Url:     captcha.PageURL,
	}
	req := cap.ToRequest()
	req.SetProxy("HTTPS", captcha.Proxy)

	return solver.SubmitRequest(req)
}

// SolveReCaptchaV2 requests a ReCaptcha v2 solution.
func (solver *CaptchaSolver) SolveReCaptchaV2(captcha *base.ReCaptchaV2) (string, error) {

	cap := api2captcha.ReCaptcha{
		SiteKey:   captcha.SiteKey,
		Url:       captcha.PageURL,
		Invisible: captcha.Invisible,
		Action:    captcha.Action,
	}
	req := cap.ToRequest()
	req.SetProxy("HTTPS", captcha.Proxy)

	return solver.SubmitRequest(req)
}

// SolveReCaptchaV2E is unsupported by AutoSolve but must be implemented to satisfy the base
// CaptchaSolver interface.
func (solver *CaptchaSolver) SolveReCaptchaV2E(captcha *base.ReCaptchaV2E) (string, error) {
	return "", ErrCaptchaUnsupported
}

// SolveReCaptchaV3E is used to get a token for the enterprise reCaptcha V3 CAPTCHA
func (solver *CaptchaSolver) SolveReCaptchaV3E(captcha *base.ReCaptchaV3E) (string, error) {
	return "", ErrCaptchaUnsupported
}

// SolveGeetest is used to get a token for either the checkbox or invisible reCaptcha V2 CAPTCHA
func (solver *CaptchaSolver) SolveGeetest(captcha *base.GeetestCaptcha) (anticaptcha.GeeTestSolution, error) {
	cap := api2captcha.GeeTest{
		GT:        captcha.GtKey,
		ApiServer: captcha.ApiServer,
		Challenge: captcha.ChallengeKey,
		Url:       captcha.PageURL,
	}
	req := cap.ToRequest()
	req.SetProxy("HTTPS", captcha.Proxy)

	code, err := solver.SubmitRequest(req)
	if err != nil {
		return anticaptcha.GeeTestSolution{}, err
	}

	var geetestSolution anticaptcha.GeeTestSolution

	// func Unmarshal(data []byte, v interface{}) error
	err = json.Unmarshal([]byte(code), &geetestSolution)
	if err != nil {
		return anticaptcha.GeeTestSolution{}, err
	}

	return geetestSolution, err
}

// SolveReCaptchaV3 requests a ReCaptcha v3 solution.
func (solver *CaptchaSolver) SolveReCaptchaV3(captcha *base.ReCaptchaV3) (string, error) {
	cap := api2captcha.ReCaptcha{
		SiteKey: captcha.SiteKey,
		Url:     captcha.PageURL,
		Score:   float64(captcha.MinScore),
		Version: "v3",
		Action:  captcha.Action,
	}
	req := cap.ToRequest()
	req.SetProxy("HTTPS", captcha.Proxy)

	return solver.SubmitRequest(req)
}
