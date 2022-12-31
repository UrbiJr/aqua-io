package twocaptcha

import "errors"

// Error related to 2Captcha API responses
var (
	ErrNoSlotAvailable    = errors.New("The queue is too long or you've exceed the maximum rate")
	ErrMaxUserTurn        = errors.New("More than 60 requests within 3 seconds. Banned for 10 seconds")
	ErrWrongIDFormat      = errors.New("The provided captcha ID is invalid. The ID can contain numbers only")
	ErrWrongCaptchaID     = errors.New("Invalid captcha ID")
	ErrBadDuplicates      = errors.New("Maximum number of attempts reached")
	ErrIPAddress          = errors.New("IP address doesn't match the pingback IP or domain")
	ErrEmptyAction        = errors.New("Action parameter is missing or no value is provided for action parameter")
	ErrPageURL            = errors.New("The page URL parameter was missing from the request")
	ErrBadTokenOrPageURL  = errors.New("The request contains invalid pair of siteKey and page URL")
	ErrGoogleKey          = errors.New("The siteKey value provided is incorrect: it's blank or malformed")
	ErrWrongUserKey       = errors.New("The provided key parameter value is invalid, it should contain 32 symbols")
	ErrKeyDoesNotExist    = errors.New("The key you've provided does not exists")
	ErrZeroBalance        = errors.New("You don't have funds on your account")
	ErrIPBanned           = errors.New("Your IP address is banned due to many frequent attempts to access the server using wrong authorization keys")
	ErrIPNotAllowed       = errors.New("The request is sent from the IP that is not on the list of your allowed IPs")
	ErrCaptchaUnsolvable  = errors.New("Captcha Unsolvable")
	ErrNotReady           = errors.New("The task is still being processed")
	ErrTimeout            = errors.New("Request Timeout")
	ErrInvalidTaskID      = errors.New("2Captcha Error #3")
	ErrInvalidSolution    = errors.New("Invalid Token")
	ErrInvalidCode        = errors.New("2Captcha Error #4")
	ErrUnknownCode        = errors.New("2Captcha Error #5")
	ErrStatusNon2xx       = errors.New("2Captcha Error #6")
	ErrCaptchaUnsupported = errors.New("captcha type unsupported")
)

// APIErrors maps error codes from 2Captcha to their respective error messages
var APIErrors = map[string]error{
	"ERROR_NO_SLOT_AVAILABLE":    ErrNoSlotAvailable,
	"MAX_USER_TURN":              ErrMaxUserTurn,
	"ERROR_WRONG_ID_FORMAT":      ErrWrongIDFormat,
	"ERROR_WRONG_CAPTCHA_ID":     ErrWrongCaptchaID,
	"ERROR_BAD_DUPLICATES":       ErrBadDuplicates,
	"ERROR_IP_ADDRES":            ErrIPAddress,
	"ERROR_EMPTY_ACTION":         ErrEmptyAction,
	"ERROR_PAGEURL":              ErrPageURL,
	"ERROR_BAD_TOKEN_OR_PAGEURL": ErrBadTokenOrPageURL,
	"ERROR_GOOGLEKEY":            ErrGoogleKey,
	"ERROR_WRONG_USER_KEY":       ErrWrongUserKey,
	"ERROR_KEY_DOES_NOT_EXIST":   ErrKeyDoesNotExist,
	"ERROR_ZERO_BALANCE":         ErrZeroBalance,
	"IP_BANNED":                  ErrIPBanned,
	"ERROR_IP_NOT_ALLOWED":       ErrIPNotAllowed,
	"ERROR_CAPTCHA_UNSOLVABLE":   ErrCaptchaUnsolvable,
	"CAPCHA_NOT_READY":           ErrNotReady,
}

// GetAPIError returns the error from the given 2Captcha response code
func GetAPIError(code string) error {
	if code == "" {
		return ErrInvalidCode
	}

	if err, ok := APIErrors[code]; ok {
		return err
	}

	return ErrUnknownCode
}
