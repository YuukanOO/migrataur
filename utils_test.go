package migrataur

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// assertInstance is a utility class used for testing only. Just drop this file somewhere
// and enjoy.
type assertInstance struct {
	t *testing.T
}

func assert(t *testing.T) *assertInstance {
	return &assertInstance{t: t}
}

func (a *assertInstance) contains(pattern, actual string) *assertInstance {
	if !strings.Contains(actual, pattern) {
		a.t.Errorf("Expected %s to contains %s", actual, pattern)
	}

	return a
}

func (a *assertInstance) true(actual bool) *assertInstance {
	if actual == false {
		a.t.Errorf("Expected to be true! Got: %t", actual)
	}

	return a
}

func (a *assertInstance) false(actual bool) *assertInstance {
	if actual == true {
		a.t.Errorf("Expected to be false! Got: %t", actual)
	}

	return a
}

func (a *assertInstance) equals(expected, actual interface{}) *assertInstance {
	if actual != expected {
		a.t.Errorf("Expected: %s, Got: %s", expected, actual)
	}

	return a
}

func (a *assertInstance) notEquals(expected, actual interface{}) *assertInstance {
	if actual == expected {
		a.t.Errorf("Should be not equals: %s, And: %s", expected, actual)
	}

	return a
}

func (a *assertInstance) notNil(actual interface{}) *assertInstance {
	if actual == nil {
		a.t.Error("Should not be nil!")
	}

	return a
}

func (a *assertInstance) nil(actual interface{}) *assertInstance {
	if actual != nil {
		a.t.Error(actual)
	}

	return a
}

func (a *assertInstance) applied(migrations []*Migration, names ...string) *assertInstance {
	lenActual, lenExpected := len(migrations), len(names)

	a.equals(lenExpected, lenActual)

	for i := 0; i < len(migrations); i++ {
		a.contains(names[i], migrations[i].Name)
	}

	return a
}

func (a *assertInstance) exists(pathes ...string) *assertInstance {
	path := filepath.Join(pathes...)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		a.t.Errorf("File %s does not exists", path)
	}

	return a
}
