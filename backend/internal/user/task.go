package user

// Trader contains information specific to a trader (copied or not)
type Task struct {
	ID          int64  `json:"id"`
	ProfileID   int64  `json:"profile_id"`
	ProxyListID int64  `json:"proxy_list_id"`
	Module      string `json:"module"`
	PaymentMode string `json:"payment_mode"`
	Running     bool   `json:"running"`
}
