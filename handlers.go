package main

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	json_header = "application/json; charset=utf-8"
)

func getQueryValue(v url.Values, param string, query *string) error {

	if qval, ok := v[param]; ok {
		if *query != "" {
			return fmt.Errorf("Too many queries %v, %v", qval, query)
		}

		if len(qval) > 1 {
			return fmt.Errorf("Only accepts a single query parameter %v", qval)
		}
		*query = qval[0]
	}

	return nil
}

func getQueryParamsAsInt(v url.Values, param string) (values []int) {

	if value, ok := v[param]; ok {
		for _, i := range value {
			v, err := strconv.Atoi(i)
			if err == nil {
				values = append(values, v)
			}
		}
	}
	return
}

func LocationSearchHandler(w http.ResponseWriter, req *http.Request) {

	hasWritten := false
	key := [2]int{}
	query := ""
	start := time.Now().UnixNano()
	v := req.URL.Query()

	w.Header().Set("Content-Type", json_header)

	postcodes := getQueryParamsAsInt(v, "postcode")
	numbers := getQueryParamsAsInt(v, "number")
	err := getQueryValue(v, "name", &query)

	if err != nil {
		log.Println(err)
		w.Write([]byte("Error"))
		return
	}

	// If there are no numbers, use the default values
	if len(numbers) == 0 {
		numbers = DefaultNumbers
	}

	w.Write([]byte("["))

	//TODO, default postcodes
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
