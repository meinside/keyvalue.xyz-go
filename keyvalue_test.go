package keyvalue

import (
	"encoding/json"
	"testing"
)

const (
	KeyForTesting = "keyvalue_xyz-key-for-testing" // XXX - put yours here
)

type obj struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
}

// $ go test

func TestKeyValue(t *testing.T) {
	val := "some-value"

	if kv, err := NewKeyValue(KeyForTesting); err == nil {
		// Set
		if err := kv.SetAndValidate(val); err != nil {
			t.Error("Failed to set value:", val, "-", err)
		}

		// Get
		if _, err := kv.Get(); err != nil {
			t.Error("Failed to get value for key:", kv.Key, "-", err)
		}
	} else {
		t.Error("Failed to generate a new key:", err)
	}
}

func TestKeyValueObj(t *testing.T) {
	val1 := obj{Name: "Tester1", Age: 37, Occupation: "Developer"}
	val2 := obj{Name: "Tester2", Age: 39, Occupation: "Entrepreneur"}

	if kv, err := NewKeyValue(KeyForTesting); err == nil {
		// Set and validate with function
		if err := kv.SetObjectAndValidate(val1, func(v string, o interface{}) bool {
			var r obj // returned object
			if err := json.Unmarshal([]byte(v), &r); err != nil {
				t.Error("Failed to unmarshal value:", v, "-", err)
				return false
			}

			if r.Name == val1.Name && r.Age == val1.Age && r.Occupation == val1.Occupation {
				return true
			}

			t.Error("Unmarshalled object is different from the requested object:", r, "-", val1)
			return false
		}); err != nil {
			t.Error(err)
		}

		// Get
		if _, err := kv.Get(); err != nil {
			t.Error(err)
		}

		// Set again (not validating)
		if err := kv.SetObject(val2); err != nil {
			t.Error(err)
		}

		// Get
		if _, err := kv.Get(); err != nil {
			t.Error(err)
		}
	} else {
		t.Error(err)
	}
}
