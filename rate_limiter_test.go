package main

import "testing"

func TestHashClientId(t *testing.T) {
	if hashClientId("client A") == hashClientId("client B") {
		t.Error("Unexpected hash collision!")
	}
}

func TestAddRateLogEntry(t *testing.T) {
	oldLen := len(rateLog)
	addRateLogEntry(0)
	if len(rateLog) != oldLen+1 {
		t.Errorf("rateLog was not appended")
	}
}

func TestCollectRateLogEntries(t *testing.T) {
	// TODO
}

func TestRateLimit(t *testing.T) {
	// TODO
}
