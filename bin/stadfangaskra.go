package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/StefanKjartansson/stadfangaskra"
	log "github.com/llimllib/loglevel"
	"net/http"
	"os"
)

func main() {

	var httpListen, logLevel string

	flag.StringVar(&httpListen, "http", "127.0.0.1:3999", "host:port to listen on")
	flag.StringVar(&logLevel, "level", "info", "log level")

	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatalf("Missing file argument")
	}

	log.SetPriorityString(logLevel)

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	locs := []stadfangaskra.Location{}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&locs)
	if err != nil {
		log.Fatal(err)
	}

	for idx, l := range locs {
		locs[idx].Name = fmt.Sprintf("%s, %d %s", l.Name, l.Postcode, l.Municipality)
		b, err := json.Marshal(l)
		if err != nil {
			panic(err)
		}
		locs[idx].JSONCache = b
	}

	log.Infof("Starting server on: %s", httpListen)
	ls := stadfangaskra.NewLocationService("/locations/", locs)
	http.Handle("/", ls.GetRouter())
	log.Fatal(http.ListenAndServe(httpListen, nil))
}
