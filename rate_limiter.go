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

func init() {
	rateLog = make([]RateLogEntry, 0)
}

func hashClientId(clientId string) uint32 {
	return crc32.ChecksumIEEE([]byte(clientId))
}

func addRateLogEntry(hash uint32) {
	rateLog = append(rateLog, RateLogEntry{time.Now().Unix(), hash})
}

// Expunge old rate log entries and return all entries for given clientId
func collectRateLogEntries(hash uint32) []RateLogEntry {
	collection := make([]RateLogEntry, 0)
	newRateLog := make([]RateLogEntry, 0)
	windowStart := time.Now().Unix() - Window

	for _, entry := range rateLog {
		if entry.timestamp > windowStart {
			newRateLog = append(newRateLog, entry)
			if entry.clientId == hash {
				collection = append(collection, entry)
			}
		}
	}

	rateLog = newRateLog

	return collection
}

func RateLimit(clientId string) bool {
	hash := hashClientId(clientId)
	clientRequests := len(collectRateLogEntries(hash))

	if clientRequests >= WindowLimit {
		log.Printf("Request for %s blocked: %d request in list", clientId, clientRequests)
		return true
	} else {
		addRateLogEntry(hash)
		return false
	}
}
