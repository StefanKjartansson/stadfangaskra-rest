package main

import (
	"code.google.com/p/gorilla/mux"
	"flag"
	"github.com/StefanKjartansson/stadfangaskra"
	"log"
	"net/http"
)

var (
	router = new(mux.Router)
)

func main() {
	var http_listen, placename_file string
	flag.StringVar(&http_listen, "http", "127.0.0.1:3999", "host:port to listen on")
	flag.StringVar(&placename_file, "file", "Stadfangaskra_20130326.dsv", "csv input file")
	flag.Parse()

	stadfangaskra.ImportDatabase(placename_file)
	log.Println("Data Imported")
	log.Println("Starting server")

	stadfangaskra.SetupRoutes(router)

	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(http_listen, nil))
	log.Println("Bye")
}
