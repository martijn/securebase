package main

import "testing"

func TestDatastore(t *testing.T) {
  datastore.Set("key", "value")
  _, str := datastore.Get("key")

  if str != "value" {
    t.Error("Expected value to be set")
  }

  datastore.Delete("key")

  _, str = datastore.Get("key")

  if str != "" {
    t.Error("Expected value to be empty")
  }

}
