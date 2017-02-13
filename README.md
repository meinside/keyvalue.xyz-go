# keyvalue.xyz-go

Go library for [keyvalue.xyz](https://keyvalue.xyz/) service.

# Install

```bash
$ go get -u github.com/meinside/keyvalue.xyz-go
```

# How-to-use

## Sample 1: Create a new key-value with string data

```go
package main

import (
	"fmt"

	"github.com/meinside/keyvalue.xyz-go"
)

func main() {
	key := "some-key-for-testing-this-library"
	val1 := "some-value"
	val2 := "other-value"

	if kv, err := keyvalue.NewKeyValue(key); err == nil {
		token := kv.Token

		fmt.Printf("> Token for key: %s = %s\n", kv.Key, token)

		// Set
		if err := kv.SetAndValidate(val1); err != nil {
			fmt.Printf("* Error while setting value: %s - %s\n", val1, err)
		}

		// Get
		if val, err := kv.Get(); err != nil {
			fmt.Printf("* Error while getting value for key: %s - %s\n", kv.Key, err)
		} else {
			fmt.Printf("> Value for key: %s = %s\n", kv.Key, val)
		}

		// Set again
		if err := kv.SetAndValidate(val2); err != nil {
			fmt.Printf("* Error while updating value to: %s - %s\n", val2, err)
		}

		// Get
		if val, err := kv.Get(); err != nil {
			fmt.Printf("* Error while getting updated value for key: %s - %s\n", kv.Key, err)
		} else {
			fmt.Printf("> Value for key: %s = %s\n", kv.Key, val)
		}
	} else {
		fmt.Printf("* Error generating a new key: %s\n", err)
	}
}
```

## Sample 2: Using a saved key-value with JSON data

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/meinside/keyvalue.xyz-go"
)

type obj struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

func main() {
	key := "some-key-for-testing-this-library"
	token := "abc01234"	// XXX - put yours here
	val1 := obj{Name: "Tester1", Age: 37, Occupation: "Software Engineer"}
	val2 := obj{Name: "Tester2", Age: 39, Occupation: "Entrepreneur"}

	kv := keyvalue.NewKeyValueWithToken(token, key)

	// Set and validate with function
	if err := kv.SetObjectAndValidateFunc(val1, func(v string, o interface{}) bool {
		var r obj // returned object
		if err := json.Unmarshal([]byte(v), &r); err != nil {
			fmt.Printf("* Error while unmarshalling value: %s - %s\n", v, err)
		}

		if r.Name == val1.Name && r.Age == val1.Age && r.Occupation == val1.Occupation {
			return true
		} else {
			fmt.Printf("* Unmarshalled object is different from requested object: %+v - %+v\n", r, val1)
			return false
		}
	}); err != nil {
		fmt.Printf("* Error while setting value: %s - %s\n", val1, err)
	}

	// Get
	if val, err := kv.Get(); err != nil {
		fmt.Printf("* Error while getting value for key: %s - %s\n", kv.Key, err)
	} else {
		fmt.Printf("> Value for key: %s = %s\n", kv.Key, val)
	}

	// Set again (not validating)
	if err := kv.SetObject(val2); err != nil {
		fmt.Printf("* Error while updating value to: %s - %s\n", val2, err)
	}

	// Get
	if val, err := kv.Get(); err != nil {
		fmt.Printf("* Error while getting updated value for key: %s - %s\n", kv.Key, err)
	} else {
		fmt.Printf("> Value for key: %s = %s\n", kv.Key, val)
	}
}
```

# Known issues / Todos

- [ ] Requests for keys with `.` return nothing (eg. "my.creative.keyname")

# License

MIT
