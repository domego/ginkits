package log

import (
	"testing"
)

func TestTracef(t *testing.T) {
	Tracef("%s testing.", "trace")
}

func TestInfof(t *testing.T) {
	Infof("%s testing.", "info")
}

func TestWarnf(t *testing.T) {
	Warnf("%s testing", "warn")
}

func TestErrorf(t *testing.T) {
	Errorf("%s testing.", "error")
}
