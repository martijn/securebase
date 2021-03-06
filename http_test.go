package main

import "bytes"
import "io/ioutil"
import "net/http"
import "net/http/httptest"
import "testing"

func TestGet(t *testing.T) {
	/* Setup */
	ts := httptest.NewServer(http.HandlerFunc(handleGet))
	defer ts.Close()

	/* Get nonexistent key */

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 404 {
		t.Error("Expected 404 response")
	}

	/* Get existing key */

	datastore.Set("gettest", encrypt("success", ""))
	res, err = http.Get(ts.URL + "/gettest")
	if err != nil {
		t.Error(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}
	if string(body) != "success" {
		t.Errorf("Unexpected result: %s", body)
	}
}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePost))
	defer ts.Close()

	res, err := http.Post(ts.URL+"/postkey", "text/plain", bytes.NewBufferString("value"))
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Expected 200 response")
	}

	_, value := datastore.Get("postkey")

	if value == "" {
		t.Error("No value found in datastore")
	}
}

func TestDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleDelete))
	defer ts.Close()

	datastore.Set("deltest", encrypt("success", ""))
	req, err := http.NewRequest(http.MethodDelete, ts.URL+"/deltest", nil)
	if err != nil {
		t.Error(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("Expected 200 response")
	}

	_, value := datastore.Get("deltest")

	if value != "" {
		t.Errorf("Value was not deleted from datastore: %x", value)
	}
}

func TestRouterRateLimit(t *testing.T) {
	for i := 0; i <= WindowLimit; i++ {
		RateLimit("127.0.0.1")
	}
	defer resetRateLog()

	ts := httptest.NewServer(http.HandlerFunc(router))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusTooManyRequests {
		t.Errorf("Expected 429 response, got %v", res.StatusCode)
	}
}

func TestRouter405(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(router))
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPatch, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 response, got %v", res.StatusCode)
	}

	// TODO: Test get/post/delete handling with mocks maybe?
}

func TestGetIP(t *testing.T) {
	if getIP("[::1]:12345") != "[::1]" {
		t.Error("Did not get expected output [::1]")
	}

	if getIP("127.0.0.1:12345") != "127.0.0.1" {
		t.Error("Did not get expected output 127.0.0.1")
	}
}
