package user

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ProxyProfile struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Proxies []string `json:"proxies"`
}

// ReadProxies reads proxies.json and returns read data as []Profile
func ReadProxies() ([]ProxyProfile, error) {

	var proxies []ProxyProfile
	path, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}
	// windows: C:\Users\<user>\AppData\Local\Roaming\NyxAIO\proxies.json
	path = filepath.Join(path, "NyxAIO", "proxies.json")
	jsonFile, err := os.Open(path)

	if err != nil {
		// file does not exist
		proxies = []ProxyProfile{}
		file, _ := json.MarshalIndent(proxies, "", " ")
		_ = ioutil.WriteFile(path, file, 0644)
		return proxies, nil
	}

	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &proxies)
	if err != nil {
		// error reading file
		return []ProxyProfile{}, err
	}

	return proxies, nil

}

// WriteProxies writes profiles to proxies.json
func WriteProxies(proxies []ProxyProfile) {
	path, err := os.UserCacheDir()
	if err != nil {
		return
	}
	// windows: C:\Users\<user>\AppData\Local\Roaming\NyxAIO\proxies.json
	path = filepath.Join(path, "NyxAIO", "proxies.json")
	file, _ := json.MarshalIndent(proxies, "", " ")
	_ = ioutil.WriteFile(path, file, 0644)
}
