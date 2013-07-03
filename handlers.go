package main

import (
	"code.google.com/p/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	json_header = "application/json; charset=utf-8"
)

func LocationSearchHandler(w http.ResponseWriter, req *http.Request) {

	hasWritten := false
	key := [2]int{}
	start := time.Now().UnixNano()
	v := req.URL.Query()

	w.Header().Set("Content-Type", json_header)

	postcodes := getQueryParamsAsInt(v, "postcode")
	numbers := getQueryParamsAsInt(v, "number")
	query, err := getSingleQueryValueOrEmpty(v, "name")

	if err != nil {
		log.Println(err)
		w.Write([]byte("Error"))
		return
	}

	// If there are no numbers, use the default values
	if len(numbers) == 0 {
		numbers = DefaultNumbers
	}

	if len(postcodes) == 0 {
		postcodes = DefaultPostCodes
	}

	w.Write([]byte("["))

	for _, pc := range postcodes {
		key[0] = pc

		for _, n := range numbers {
			key[1] = n

			for _, l := range LookupTable[key] {

				if !l.NameMatches(query) {
					continue
				}

				if hasWritten {
					w.Write([]byte(","))
				}

				w.Write(l.JSONCache)
				hasWritten = true

			}
		}
	}

	w.Write([]byte("]"))

	log.Printf("%s %s %s, time: %f.ms", req.RemoteAddr,
		req.Method, req.URL.Query(),
		float64(time.Now().UnixNano()-start)/1000000.0)

	return
}

func LocationDetailHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", json_header)

	start := time.Now().UnixNano()
	w.Write(IndexTable[id].JSONCache)

	log.Printf("%s %s %s, time: %f.ms", req.RemoteAddr,
		req.Method, req.URL.Query(),
		float64(time.Now().UnixNano()-start)/1000000.0)

	return
}
