package migrataur

import (
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

func (a *assertInstance) migrationsEquals(migrations []*Migration, names ...string) *assertInstance {
	lenActual, lenExpected := len(migrations), len(names)

	a.equals(lenExpected, lenActual)

	for i := 0; i < len(migrations); i++ {
		a.contains(names[i], migrations[i].Name)
	}

	return a
}
