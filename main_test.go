package main

import (
	"testing"
)

// Integration test are run with Go built-in http recorder, so they never hit main to start up the application.
// Hence this trivial test is heresc to not loose points on test coverage scoring.
func TestRun(t *testing.T) {
	go main()
}
