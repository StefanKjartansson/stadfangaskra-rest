package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

func getQueryParamsAsInt(v url.Values, param string) (values []int64) {

	if value, ok := v[param]; ok {
		for _, i := range value {
			v, err := strconv.ParseInt(i, 0, 64)
			if err == nil {
				values = append(values, v)
			}
		}
	}

	return
}

func ParseQueryParams(v url.Values) (postcodes []int64, numbers []int64, query string, err error) {

	postcodes = getQueryParamsAsInt(v, "postcode")
	numbers = getQueryParamsAsInt(v, "number")

	ptr := &query
	err = getQueryValue(v, "name_startswith", ptr)
	if err != nil {
		return
	}

	err = getQueryValue(v, "name_endswith", ptr)
	if err != nil {
		return
	}

	err = getQueryValue(v, "name", ptr)
	if err != nil {
		return
	}
	return
}

func GetLocation(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL.Query())
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

	for _, element := range Locations {
		if element.ContainsPostcode(postcodes) &&
			element.ContainsNumbers(numbers) &&
			element.NameMatches(query) {
			if hasWritten {
				w.Write([]byte(","))
			}
			if err := enc.Encode(&element); err != nil {
				log.Println(err)
			}
			hasWritten = true
		}
	}
	w.Write([]byte("]"))
	return
}
