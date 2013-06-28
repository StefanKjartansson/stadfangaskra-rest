package main

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	nameQueries = []string{"name_startswith", "name_endswith", "name"}
)

func getQueryValue(v url.Values, param string, query *string) error {

	if qval, ok := v[param]; ok {
		if *query != "" {
			return fmt.Errorf("Too many queries %v, %v", qval, query)
		}

		if len(qval) > 1 {
			return fmt.Errorf("Only accepts a single query parameter %v", qval)
		}
		if strings.HasSuffix(param, "_startswith") {
			*query = qval[0] + "*"
		} else if strings.HasSuffix(param, "_endswith") {
			*query = "*" + qval[0]
		} else {
			*query = qval[0]
		}
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

func ParseQueryParams(v url.Values) (postcodes []int, numbers []int, query string, err error) {

	postcodes = getQueryParamsAsInt(v, "postcode")
	numbers = getQueryParamsAsInt(v, "number")

	for _, i := range nameQueries {
		err = getQueryValue(v, i, &query)
		if err != nil {
			return
		}
	}
	return
}

func LocationSearchHandler(w http.ResponseWriter, req *http.Request) {

	start := time.Now().UnixNano()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	postcodes, numbers, query, err := ParseQueryParams(req.URL.Query())

	if err != nil {
		log.Println(err)
		w.Write([]byte("Error"))
		return
	}

	w.Write([]byte("["))

	enc := json.NewEncoder(w)
	hasWritten := false

	for _, i := range Locations {
		if i.ContainsPostcode(postcodes) &&
			i.ContainsNumbers(numbers) &&
			i.NameMatches(query) {

			if hasWritten {
				w.Write([]byte(","))
			}
			if err := enc.Encode(&i); err != nil {
				log.Println(err)
			}
			hasWritten = true
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
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	for _, i := range Locations {
		if i.Hnitnum == id {
			b, err := json.Marshal(i)
			if err != nil {
				fmt.Println("error:", err)
			}
			w.Write(b)
			return
		}
	}
	w.Write([]byte("{}"))
	return
}
