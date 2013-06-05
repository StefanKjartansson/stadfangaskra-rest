package main

import (
	"code.google.com/p/gorilla/mux"
	"flag"
	"log"
	"net/http"
)

var (
	httpListen    = flag.String("http", "127.0.0.1:3999", "host:port to listen on")
	placenameFile = flag.String("file", "Stadfangaskra_20130326.dsv", "csv input file")
	Locations     []Location
)

func main() {
	flag.Parse()
	log.Println("Starting import")
	ImportDatabase(*placenameFile)
	log.Println("Data Imported")
	log.Println("Starting server")

	r := mux.NewRouter()
	r.HandleFunc("/locations/",
		LocationSearchHandler)
	r.HandleFunc("/locations/{id:[0-9]+}/",
		LocationDetailHandler)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(*httpListen, nil))
	log.Println("Bye")
}
