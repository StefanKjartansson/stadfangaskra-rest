package stadfangaskra

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func timedResponse(f func(http.ResponseWriter, *http.Request)) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

		start := time.Now().UnixNano()
		f(w, r)
		log.Printf("%s %s %s, time: %f.ms", r.RemoteAddr,
			r.Method, r.URL.Query(),
			float64(time.Now().UnixNano()-start)/1000000.0)
	}
}

func errorHandler(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("handling %q: %v", r.RequestURI, err)
		}
	}
}

func byteChanToJSON(w io.Writer, cs chan []byte, quoteValue bool) chan bool {

	done := make(chan bool, 1)

	go func() {

		quote := []byte("\"")
		hasWritten := false
		w.Write([]byte("["))
		for s := range cs {
			if hasWritten {
				w.Write([]byte(","))
			}
			if quoteValue {
				w.Write(quote)
			}
			w.Write(s)
			if quoteValue {
				w.Write(quote)
			}
			hasWritten = true
		}
		w.Write([]byte("]"))

		done <- true
	}()

	return done
}

func LocationSearchHandler(w http.ResponseWriter, req *http.Request) error {

	key := [2]int{}
	v := req.URL.Query()

	postcodes := getQueryParamsAsInt(v, "postcode")
	numbers := getQueryParamsAsInt(v, "number")
	query, err := getSingleQueryValueOrEmpty(v, "name")

	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	// If there are no numbers, use the default values
	if len(numbers) == 0 {
		numbers = DefaultNumbers
	}

	if len(postcodes) == 0 {
		postcodes = DefaultPostCodes
	}

	cs := make(chan []byte)
	done := byteChanToJSON(w, cs, false)
	for _, pc := range postcodes {
		key[0] = pc
		for _, n := range numbers {
			key[1] = n
			for _, l := range LookupTable[key] {
				if !l.NameMatches(query) {
					continue
				}
				cs <- l.JSONCache
			}
		}
	}
	close(cs)
	<-done
	return nil

}

func LocationDetailHandler(w http.ResponseWriter, req *http.Request) error {

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return fmt.Errorf("Invalid id parameter")
	}
	w.Write(IndexTable[id].JSONCache)
	return nil
}

func AutoCompleteStreetNamesHandler(w http.ResponseWriter, req *http.Request) error {

	v := req.URL.Query()
	q, present := v["q"]
	if !present {
		return fmt.Errorf("Query parameter missing")
	}

	matcher := q[0]

	cs := make(chan []byte)
	done := byteChanToJSON(w, cs, true)
	for _, s := range StreetNames {
		if strings.Contains(s, matcher) {
			cs <- []byte(s)
		}
	}
	close(cs)
	<-done
	return nil
}

func SetupRoutes(r *mux.Router) {

	r.HandleFunc("/locations/",
		timedResponse(errorHandler(LocationSearchHandler)))
	r.HandleFunc("/locations/{id:[0-9]+}/",
		timedResponse(errorHandler(LocationDetailHandler)))
	r.HandleFunc("/ac/streets/",
		timedResponse(errorHandler(AutoCompleteStreetNamesHandler)))
}
