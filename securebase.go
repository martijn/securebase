package main

var datastore *BoltDatastore

func main() {
	datastore = NewBoltDatastore("datastore")
	defer datastore.Close()

	startHttpServer()
}
