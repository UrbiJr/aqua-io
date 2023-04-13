package task

import (
	"encoding/base64"

	"github.com/UrbiJr/copy-io/internal/captcha/base"
)

// NormalCaptchaTaskParams defines parameters used to solve Normal Captchas. Normal Captcha is an image that contains distored but human-readable text.
// See https://2captcha.com/2captcha-api#solving_normal_captcha for more information.
type NormalCaptchaTaskParams struct {
	InitParams

	// Methods are `base64` and `post`
	// We only have support for the `base64` method for now.
	Method string `json:"method"`
	// Base64 encoded image data (single image)(required)
	Body string `json:"body"`
	// i.e. Click all of the sports equipment
	// If the captcha image contains instructions, this option may be omitted.
	// Size limit of 140 bytes, UTF-8
	InstructionsText string `json:"textinstructions,omitempty"`
	// This option and text instructions are mutually exclusive.
	// Size limit of 100 KiB
	InstructionsBody string `json:"imginstructions,omitempty"`
	// https://2captcha.com/2captcha-api#language
	// Optional language code i.e. `en`, `ru`, `es` ...
	Language string `json:"lang,omitempty"`
	// `false` - Captcha contains a single word (default)
	// `true` - Captcha contains multiple words
	Phrase bool `json:"phrase"`
	// `false` - Captcha is case insensitive (default)
	// `true` - Captcha is case sensitive
	CaseSensitive bool `json:"regsense"`
	// `0` - Not specified (default)
	// `1` - Captcha contains only numbers
	// `2` - Captcha contains only letters
	// `3` - Captcha contains only numbers OR only letters
	// `4` - Captcha contains both numbers AND letters
	Numeric int `json:"numeric"`
	// `false` - Not specified (default)
	// `true` - Captcha requires calculation i.e. 2 + 2
	Math bool `json:"calc"`

	// The remaining options are unused/untested

	// The filename of a captcha image file in a MIME multipart entity body.
	// This option is used with the `post` method when using the HHTP POST verb.
	// File  string `json:"file"`

	// Base64 encoded image of instructions on how to solve the captcha
	// ImageInstructions string `json:"imginstructions"`

	// `0` - Not specified (default)
	// `1` - 1..20 - minimum number of symbols in the captcha
	// MinLength int `json:"min_len"`

	// `0` - Not specified (default)
	// `1` - 1..20 - maximum number of symbols in the captcha
	// MaxLength int `json:"max_len"`
}

// NormalCaptchaTask contains data needed to perform Normal Captcha tasks.
type NormalCaptchaTask struct {
	*BaseTask
}

// NewNormalCaptchaTaskParams returns a new instance of NormalCaptchaTask.
func NewNormalCaptchaTaskParams(captcha *base.NormalCaptcha, defaults *DefaultParams) *NormalCaptchaTaskParams {
	initParams := &NormalCaptchaTaskParams{}
	initParams.SetDefaults(defaults)

	initParams.Method = "base64"
	initParams.Body = base64.StdEncoding.EncodeToString(captcha.Body)

	switch instructions := captcha.Instructions.(type) {
	case string:
		initParams.InstructionsText = instructions
	case []byte:
		initParams.InstructionsBody = base64.StdEncoding.EncodeToString(instructions)
	}

	initParams.Language = captcha.Language
	initParams.Phrase = captcha.Phrase
	initParams.CaseSensitive = captcha.CaseSensitive
	initParams.Math = captcha.Math

	if captcha.Numeric {
		initParams.Numeric = 1
	} else if captcha.Letters {
		initParams.Numeric = 2
	} else if captcha.AlphaNumericXOR {
		initParams.Numeric = 3
	} else if captcha.AlphaNumericAND {
		initParams.Numeric = 4
	}

	return initParams
}
