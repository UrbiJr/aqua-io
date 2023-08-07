package ws

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gorilla/websocket"
	"time"
)

const (
	private = "wss://socket.okx.com:8443/socket/v5/public"
	public  = "wss://socket.okx.com:8443/socket/v5/private"
)

const (
	redialTick = 2 * time.Second
	pongWait   = 30 * time.Second
	PingPeriod = (pongWait * 8) / 10
)

var (
	unableToConnect = errors.New("ERROR Unable to Establish Websocket Connection")
)

type OKXClient struct {
	Context     context.Context
	Cancel      context.CancelFunc
	Connections map[bool]*websocket.Conn
}

type Private struct {
}

func NewClient(ctx context.Context, publicAPI, privateAPI, passphrase string) *OKXClient {
	ctx, cancel := context.WithCancel(ctx)
	client := &OKXClient{
		Context: ctx,
		Cancel:  cancel,
	}

	err := client.Connect()
	if err != nil {
		log.Error(fmt.Sprintf("ERROR Establishing OKX Websocket Connection to [%s]", publicAPI))
		return nil
	}

	return client
}

func (s *OKXClient) Connect() error {
	var buf bytes.Buffer

	var loginURI = "/users/self/verify"

	var args []map[string]interface{}

	payload := map[string]interface{}{
		"op":   "login",
		"args": args,
	}

	for _, okx := range config.GlobalConfig.OKX {
		timestamp, signature := s.sign(http.MethodGet, loginURI)
		arg := map[string]interface{}{
			"apiKey":     okx.PublicAPI,
			"passphrase": okx.Passphrase,
			"timestamp":  timestamp,
			"sign":       signature,
		}
		args = append(args, arg)
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return err
	}

	conn, resp, err := websocket.DefaultDialer.Dial(private+loginURI, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != 101 {
		return unableToConnect
	}

	err = resp.Body.Close()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(redialTick)
	defer ticker.Stop()

	/*
		go func() {
			for {
				select {
				case <-ticker.C:
					err = c.dial(p)
					if err == nil {
						return nil
					}
				case <-c.ctx.Done():
					return c.handleCancel("connect")
				}
			}
		}()
	*/

	return nil
}

func (s *OKXClient) Subscribe() {

}

func (s *OKXClient) Tickers() {

	var buf bytes.Buffer

	payload := map[string]interface{}{
		"op": "subscribe",
		"args": []map[string]interface{}{
			{
				"channel": "tickers",
				"instId":  "",
			},
		},
	}

	err := json.NewEncoder(&buf).Encode(payload)
	if err != nil {
		return
	}
}

func (s *OKXClient) sign(method, path string) (string, string) {
	timeNow := time.Now().UTC().Unix()
	timestamp := fmt.Sprint(timeNow)
	sum := timestamp + method + path
	mac := hmac.New(sha256.New, []byte(s.SecretAPI))
	mac.Write([]byte(sum))
	return timestamp, base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
