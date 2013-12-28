package main

import (
	//"flag"
	"github.com/StefanKjartansson/stadfangaskra-rest/rest"
	log "github.com/llimllib/loglevel"
	"net/http"
)

func main() {

	httpListen := "127.0.0.1:3999"

	log.SetPriorityString("info")
	ls := rest.NewLocationService("/locations/")
	http.Handle("/", ls.GetRouter())

	log.Infof("Starting server on: %s", httpListen)
	log.Fatal(http.ListenAndServe(httpListen, nil))
}
