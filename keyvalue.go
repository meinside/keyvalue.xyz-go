package keyvalue

// Service provided by: https://keyvalue.xyz/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type KeyValue struct {
	Token string `json:"token"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Get a new key-value with given key string
func NewKeyValue(key string) (kv *KeyValue, err error) {
	// request format: "https://api.keyvalue.xyz/new/some-key"
	if result, err := postUrl(fmt.Sprintf("https://api.keyvalue.xyz/new/%s", url.QueryEscape(key))); err == nil {
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
	} else {
		return nil, err
	}
}

// Get a key-value with stored token and key string
func NewKeyValueWithToken(token, key string) *KeyValue {
	return &KeyValue{
		Token: token,
		Key:   key,
	}
}

// Set a value
func (kv *KeyValue) Set(value string) error {
	// request format: "https://api.keyvalue.xyz/0123456/some-key/some-value"
	if _, err := postUrl(fmt.Sprintf("https://api.keyvalue.xyz/%s/%s/%s", kv.Token, url.QueryEscape(kv.Key), url.QueryEscape(value))); err != nil {
		return err
	}

	return nil
}

// Set a value (and validate it by comparing with the changed value)
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

// Set an object as a value
func (kv *KeyValue) SetObject(obj interface{}) error {
	if bytes, err := json.Marshal(obj); err == nil {
		return kv.Set(string(bytes))
	} else {
		return err
	}
}

// Set an object (and validate it by calling the given function)
func (kv *KeyValue) SetObjectAndValidateFunc(obj interface{}, equals func(returned string, requested interface{}) bool) error {
	if bytes, err := json.Marshal(obj); err == nil {
		return kv.Set(string(bytes))
	} else {
		return err
	}

	if v, err := kv.Get(); err != nil {
		return err
	} else {
		if !equals(v, obj) {
			return fmt.Errorf("Returned object is different from the request: %s - %+v", v, obj)
		}
	}

	return nil
}

// Get a stored value
func (kv *KeyValue) Get() (string, error) {
	// request format: "https://api.keyvalue.xyz/0123456/some-key"
	if result, err := getUrl(fmt.Sprintf("https://api.keyvalue.xyz/%s/%s", kv.Token, url.QueryEscape(kv.Key))); err != nil {
		return "", err
	} else {
		return result, nil
	}
}

func postUrl(url string) (string, error) {
	return request("POST", url)
}

func getUrl(url string) (string, error) {
	return request("GET", url)
}

func request(method, url string) (result string, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, nil); err == nil {
		var resp *http.Response
		client := &http.Client{}
		if resp, err = client.Do(req); err == nil {
			defer resp.Body.Close()

			if resp.StatusCode == 200 {
				res, _ := ioutil.ReadAll(resp.Body)

				return strings.TrimSuffix(string(res), "\n"), nil
			} else {
				return "", fmt.Errorf("Request error: %s", resp.Status)
			}
		}
	}

	return "", err
}
