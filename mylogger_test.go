package myLog

import "testing"

func TestConstLevel(t *testing.T) {
	t.Logf("%v %T\n", DEBUG, DEBUG)
	t.Logf("%v %T\n", FATAL, FATAL)
}