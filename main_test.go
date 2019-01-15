package main

import "os"
import "testing"

func TestMain(m *testing.M) {
	server_secret = []byte("server-secret")
	datastore = NewBoltDatastore("test_datastore")

	code := m.Run()

	os.Remove("test_datastore")
	os.Exit(code)
}
