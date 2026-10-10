package CRUDTEMPLATE

import (
	"testing"
)

type coverageComponent struct {
	name string
}

type secondCoverageComponent struct {
	value int
}

type coverageEvent struct {
	value string
}

func TestMyFunction(t *testing.T) {
	if got := MyFunction(); got != 1 {
		t.Fatalf("MyFunction returned %#v, want 1", got)
	}
}
