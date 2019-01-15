package main

import "fmt"
import "io/ioutil"
import "log"
import "net/http"

type requestInfo struct {
	key, clientSecret string
}

func router(res http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path, req.RemoteAddr)

	switch req.Method {
	case "GET":
		handleGet(res, req)
	case "POST":
		handlePost(res, req)
	case "DELETE":
		handleDelete(res, req)
	default:
		http.Error(res, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func handleGet(res http.ResponseWriter, req *http.Request) {
	info := parseRequest(req)
	err, value := datastore.Get(info.key)
	if err != nil {
		http.Error(res, "Unexpected error while reading from datastore", http.StatusInternalServerError)
		panic(err)
	}

	if value == "" {
		http.Error(res, "Key not found", http.StatusNotFound)
	} else {
		err, value = decrypt(value, info.clientSecret)
		if err != nil {
			http.Error(res, "Error in decryption. Incorrect Client-secret?", http.StatusUnauthorized)
		} else {
			fmt.Fprintf(res, "%s", value)
		}
	}
}

func handlePost(res http.ResponseWriter, req *http.Request) {
	info := parseRequest(req)
	value, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Unexpected error while parsing request", http.StatusInternalServerError)
		panic(err)
	}

	err = datastore.Set(info.key, encrypt(string(value), info.clientSecret))
	if err != nil {
		http.Error(res, "Unexpected error while writing to datastore", http.StatusInternalServerError)
		panic(err)
	}
	res.WriteHeader(http.StatusOK)
}

func handleDelete(res http.ResponseWriter, req *http.Request) {
	info := parseRequest(req)

	err := datastore.Delete(info.key)
	if err != nil {
		http.Error(res, "Unexpected error while writing to datastore", http.StatusInternalServerError)
		panic(err)
	}
	res.WriteHeader(http.StatusOK)
}

func parseRequest(req *http.Request) *requestInfo {
	return &requestInfo{key: req.URL.Path[1:], clientSecret: req.Header.Get("Client-secret")}
}

func startHttpServer() {
	log.Println("Server starting on :5800")
	http.HandleFunc("/", router)
	log.Fatal(http.ListenAndServe(":5800", nil))
}
