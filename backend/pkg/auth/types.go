package auth

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"net/url"
)

const baseAuthUrl = ""

var (
	validate = &url.URL{
		Scheme: "https://",
		Host:   baseAuthUrl,
		Path:   "/validate",
	}
	reset = &url.URL{
		Scheme: "https://",
		Host:   baseAuthUrl,
		Path:   "/reset",
	}
)

type AuthenticationClient struct {
	client tls_client.HttpClient
}

var GlobalUserInformation UserInformation

type UserInformation struct {
	DiscordID        string
	DiscordUsername  string
	DiscordUserImage string
	ExpirationDate   string
}

type Response struct {
	Token             string `json:"token"`
	Id                string `json:"id"`
	Product           string `json:"product"`
	User              string `json:"user"`
	Plan              string `json:"plan"`
	Email             string `json:"email"`
	Status            string `json:"status"`
	Valid             bool   `json:"valid"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
	PaymentProcessor  string `json:"payment_processor"`
	LicenseKey        string `json:"license_key"`
	Metadata          struct {
		Hwid string `json:"hwid"`
	} `json:"metadata"`
	Quantity              int         `json:"quantity"`
	WalletAddress         interface{} `json:"wallet_address"`
	CustomFieldsResponses struct {
	} `json:"custom_fields_responses"`
	Discord struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		ImageUrl string `json:"image_url"`
	} `json:"discord"`
	NftTokens          []interface{} `json:"nft_tokens"`
	ExpiresAt          int64         `json:"expires_at"`
	RenewalPeriodStart interface{}   `json:"renewal_period_start"`
	RenewalPeriodEnd   interface{}   `json:"renewal_period_end"`
	CreatedAt          int           `json:"created_at"`
	ManageUrl          string        `json:"manage_url"`
	AffiliatePageUrl   string        `json:"affiliate_page_url"`
	CheckoutSession    interface{}   `json:"checkout_session"`
	AccessPass         string        `json:"access_pass"`
}

type AuthenticationResponse struct {
	Token             string `json:"token"`
	Id                string `json:"id"`
	Product           string `json:"product"`
	User              string `json:"user"`
	Plan              string `json:"plan"`
	Email             string `json:"email"`
	Status            string `json:"status"`
	Valid             bool   `json:"valid"`
	CancelAtPeriodEnd bool   `json:"cancel_at_period_end"`
	PaymentProcessor  string `json:"payment_processor"`
	LicenseKey        string `json:"license_key"`
	Metadata          struct {
		Hwid string `json:"hwid"`
	} `json:"metadata"`
	Quantity              int         `json:"quantity"`
	WalletAddress         interface{} `json:"wallet_address"`
	CustomFieldsResponses struct {
	} `json:"custom_fields_responses"`
	Discord struct {
		Id       string `json:"id"`
		Username string `json:"username"`
		ImageUrl string `json:"image_url"`
	} `json:"discord"`
	NftTokens          []interface{} `json:"nft_tokens"`
	ExpiresAt          int64         `json:"expires_at"`
	RenewalPeriodStart interface{}   `json:"renewal_period_start"`
	RenewalPeriodEnd   interface{}   `json:"renewal_period_end"`
	CreatedAt          int           `json:"created_at"`
	ManageUrl          string        `json:"manage_url"`
	AffiliatePageUrl   string        `json:"affiliate_page_url"`
	CheckoutSession    interface{}   `json:"checkout_session"`
	AccessPass         string        `json:"access_pass"`
}
