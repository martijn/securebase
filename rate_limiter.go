package main

import "log"
import "hash/crc32"
import "time"

/* A simple sliding window rate limiter */

// limit to 10 requests per minute
const Window = 60
const WindowLimit = 10

type RateLogEntry struct {
	timestamp int64
	clientId  uint32
}

var rateLog []RateLogEntry
var housekeepingTicker *time.Ticker

func init() {
	resetRateLog()
	initHousekeepingTicker()
}

func resetRateLog() {
	rateLog = make([]RateLogEntry, 0)
}

func initHousekeepingTicker() {
	housekeepingTicker = time.NewTicker(time.Second * Window * 2)
}

func hashClientId(clientId string) uint32 {
	return crc32.ChecksumIEEE([]byte(clientId))
}

func addRateLogEntry(hash uint32, timestamp int64) {
	rateLog = append(rateLog, RateLogEntry{timestamp, hash})
}

// Expunge old rate log entries and return all entries for given clientId
func collectRateLogEntries(hash uint32) []RateLogEntry {
	collection := make([]RateLogEntry, 0)
	newRateLog := make([]RateLogEntry, 0)
	windowStart := time.Now().Unix() - Window
	var housekeeping bool

	select {
	case _ = <-housekeepingTicker.C:
		log.Printf("Starting rate log housekeeping. Current size: %d", len(rateLog))
		housekeeping = true
	default:
		housekeeping = false
	}

	for _, entry := range rateLog {
		if entry.timestamp >= windowStart {
			if housekeeping {
				newRateLog = append(newRateLog, entry)
			}
			if entry.clientId == hash {
				collection = append(collection, entry)
			}
		}
	}

	if housekeeping {
		rateLog = newRateLog
		log.Printf("Finished rate log housekeeping. Current size: %d", len(rateLog))
	}

	return collection
}

// Returns true if limit is hit, false if request is allowed
func RateLimit(clientId string) bool {
	hash := hashClientId(clientId)
	clientRequests := len(collectRateLogEntries(hash))

	if clientRequests >= WindowLimit {
		log.Printf("Request for %s blocked: %d request in list", clientId, clientRequests)
		return true
	} else {
		addRateLogEntry(hash, time.Now().Unix())
		return false
	}
}
