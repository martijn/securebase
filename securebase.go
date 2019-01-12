package main

var datastore *BoltDatastore

func main() {
	datastore = NewBoltDatastore()
	defer datastore.Close()

	startHttpServer()
}
