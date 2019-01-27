package main

import "testing"
import "time"

func TestHashClientId(t *testing.T) {
	if hashClientId("client A") == hashClientId("client B") {
		t.Error("Unexpected hash collision!")
	}
}

func TestAddRateLogEntry(t *testing.T) {
	oldLen := len(rateLog)
	addRateLogEntry(0, time.Now().Unix())
	if len(rateLog) != oldLen+1 {
		t.Errorf("rateLog was not appended")
	}
}

func TestCollectRateLogEntries(t *testing.T) {
	resetRateLog()

	addRateLogEntry(0, time.Now().Unix())
	addRateLogEntry(0, time.Now().Unix()-Window-1)

	// Test collection

	if len(rateLog) != 2 {
		t.Errorf("rateLog len was %d, expected 2", len(rateLog))
	}
	if len(collectRateLogEntries(0)) != 1 {
		t.Errorf("collection len was %d, expected 1", len(collectRateLogEntries(0)))
	}

	// Test housekeeping

	// Speed up timer, rest to default behaviour afterward
	housekeepingTicker.Stop()
	housekeepingTicker = time.NewTicker(time.Nanosecond)
	defer housekeepingTicker.Stop()
	defer initHousekeepingTicker()
	time.Sleep(time.Millisecond)

	if len(collectRateLogEntries(0)) != 1 {
		t.Errorf("collection len was %d, expected 1", len(collectRateLogEntries(0)))
	}
	if len(rateLog) != 1 {
		t.Errorf("rateLog len was %d, expected 1", len(rateLog))
	}
}

func TestRateLimit(t *testing.T) {
	resetRateLog()

	for i := 0; i <= WindowLimit-2; i++ {
		RateLimit("")
	}

	if RateLimit("") {
		t.Error("Expected rate limiter to not trigger")
	}

	if !RateLimit("") {
		t.Error("Expected rate limiter to trigger")
	}

}
