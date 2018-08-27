package keyvalue

// Service provided by: https://keyvalue.xyz/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// KeyValue struct
type KeyValue struct {
	Token string `json:"token"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewKeyValue gets a new key-value with given key string
func NewKeyValue(key string) (kv *KeyValue, err error) {
	// request format: "https://api.keyvalue.xyz/new/some-key"
	result, err := postURL(fmt.Sprintf("https://api.keyvalue.xyz/new/%s", url.QueryEscape(key)))
	if err == nil {
		// result format: "https://api.keyvalue.xyz/0123456/some-key"
		splits := strings.Split(result, "/")
		cnt := len(splits)

		// sanity check 1
		if cnt <= 2 {
			return nil, fmt.Errorf("Wrong response from server: %s", result)
		}
		t, k := splits[cnt-2], splits[cnt-1] // token: "0123456", key: "some-key"

		// sanity check 2
		if k != key {
			return nil, fmt.Errorf("Returned key is different from the request: %s - %s", k, key)
		}

		return &KeyValue{
			Token: t,
			Key:   k,
		}, nil
	}

	return nil, err
}

// NewKeyValueWithToken gets a key-value with stored token and key string
func NewKeyValueWithToken(token, key string) *KeyValue {
	return &KeyValue{
		Token: token,
		Key:   key,
	}
}

// Set sets a value
func (kv *KeyValue) Set(value string) error {
	// request format: "https://api.keyvalue.xyz/0123456/some-key/some-value"
	if _, err := postURL(fmt.Sprintf("https://api.keyvalue.xyz/%s/%s/%s", kv.Token, url.QueryEscape(kv.Key), url.QueryEscape(value))); err != nil {
		return err
	}

	return nil
}

// SetAndValidate sets a value (and validate it by comparing with the changed value)
func (kv *KeyValue) SetAndValidate(value string) error {
	if err := kv.Set(value); err != nil {
		return err
	}

	if v, err := kv.Get(); err != nil {
		return err
	} else if v != value {
		return fmt.Errorf("Returned value is different from the request: %s - %s", v, value)
	}

	return nil
}

// SetObject sets an object as a value
func (kv *KeyValue) SetObject(obj interface{}) error {
	bytes, err := json.Marshal(obj)
	if err == nil {
		return kv.Set(string(bytes))
	}

	return err
}

// SetObjectAndValidate sets an object and validate it by calling the given function
func (kv *KeyValue) SetObjectAndValidate(obj interface{}, equals func(returned string, requested interface{}) bool) error {
	var err error
	var bytes []byte

	bytes, err = json.Marshal(obj)
	if err == nil {
		err = kv.Set(string(bytes))
		if err == nil {
			var val string
			val, err = kv.Get()
			if err == nil {
				if equals(val, obj) {
					return nil
				}

				return fmt.Errorf("Validation failed: `%s` differs from `%+v`", val, obj)
			}
		}
	}

	return err
}

// Get gets a stored value
func (kv *KeyValue) Get() (string, error) {
	// request format: "https://api.keyvalue.xyz/0123456/some-key"
	result, err := getURL(fmt.Sprintf("https://api.keyvalue.xyz/%s/%s", kv.Token, url.QueryEscape(kv.Key)))

	if err != nil {
		return "", err
	}

	return result, nil
}

func postURL(url string) (string, error) {
	return request("POST", url)
}

func getURL(url string) (string, error) {
	return request("GET", url)
}

func request(method, url string) (result string, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, nil); err == nil {
		client := &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 300 * time.Second,
				}).Dial,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}

		var resp *http.Response
		resp, err = client.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err == nil {
			if resp.StatusCode == 200 {
				res, _ := ioutil.ReadAll(resp.Body)

				return strings.TrimSuffix(string(res), "\n"), nil
			}

			return "", fmt.Errorf("Request error: %s", resp.Status)
		}
	}

	return "", err
}
