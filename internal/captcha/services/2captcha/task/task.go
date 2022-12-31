package task

// Task interface defines methods used by all Task type objects.
type Task interface {
	GetInitParams() Params
	GetInitResult() *Result
	SetID(id int)
	GetID() int
	GetPollParams() Params
	SetPollResult(result *Result)
	GetPollResult() *Result
	GetSolution() (string, error)
}

// Params interface defines methods used by all Params objects.
type Params interface {
	SetDefaults(defaults *DefaultParams)
}

// DefaultParams defines default parameters that must be included in solve requests
type DefaultParams struct {
	// User API key
	APIKey string `json:"key" schema:"key"`
	// `0` - The response encoding will be text/plain.
	// `1` - The response encoding will be application/json.
	JSON int `json:"json" schema:"json"`
	// If enabled, the server will send back the params that it received.
	// https://2captcha.com/2captcha-api#debugging
	// `0` - Disabled, `1` - Enabled
	Debug int `json:"debug_dump" schema:"-"`

	// If enabled `res.php` will include Access-Control-Allow-Origin:* header in the response.
	// Used for cross-domain AJAX requests in web applications.
	// `0` - Disabled, `1` - Enabled
	// ACAO	int `json:"header_acao"`
}

// SetDefaults sets default parameters used in every task.
func (params *DefaultParams) SetDefaults(defaults *DefaultParams) {
	params.APIKey = defaults.APIKey
	params.Debug = defaults.Debug
	params.JSON = 1
}

// InitParams extends TaskParams with additional optional parameters.
type InitParams struct {
	DefaultParams
	// The identifier of this software.
	SoftID int `json:"soft_id"`

	// URL for pingback (callback) response that will be sent when captcha is solved.
	// URL should be registered on the server.
	// PingBack string `json:"pingback"`

}

// SetDefaults sets the default parameters to be used with every task.
func (params *InitParams) SetDefaults(defaults *DefaultParams) {
	params.DefaultParams.SetDefaults(defaults)
	params.SoftID = 2562
}

// PollParams defines parameters to poll for a task solution.
type PollParams struct {
	DefaultParams
	// The task ID
	ID int `schema:"id"`
	// The Task action: `get`, `reportgood`, or `reportbad`
	Action string `schema:"action"`
}

// SetID sets task ID
func (params *PollParams) SetID(id int) {
	params.ID = id
}

// Result holds the response from polling a task.
type Result struct {
	Status       int    `json:"status"`
	Value        string `json:"request"`
	ErrorMessage string `json:"error_text"`
}

// BaseTask object contains data and methods needed to initiate and poll for a 2Captcha solution.
type BaseTask struct {
	ID         int
	InitParams Params
	PollParams interface {
		Params
		SetID(id int)
	}
	InitResult Result
	PollResult *Result
}

// GetInitParams returns the parameters used to initiate a solve request.
func (vt *BaseTask) GetInitParams() Params {
	return vt.InitParams
}

// GetInitResult returns the response given by 2Captcha when initiating a solve request.
func (vt *BaseTask) GetInitResult() *Result {
	return &vt.InitResult
}

// GetPollParams returns the parameters used to poll for the task solution.
func (vt *BaseTask) GetPollParams() Params {
	vt.PollParams.SetID(vt.ID)
	return vt.PollParams
}

// SetPollResult sets PollResult with the response given when polling for the task solution.
func (vt *BaseTask) SetPollResult(result *Result) {
	vt.PollResult = result
}

// GetPollResult returns the response of the last poll request.
func (vt *BaseTask) GetPollResult() *Result {
	return vt.PollResult
}

// SetID sets ID.
func (vt *BaseTask) SetID(id int) {
	vt.ID = id
}

// GetID returns ID.
func (vt *BaseTask) GetID() int {
	return vt.ID
}

// GetSolution returns the last poll result/solution for the task.
func (vt *BaseTask) GetSolution() (string, error) {
	return vt.PollResult.Value, nil
}
