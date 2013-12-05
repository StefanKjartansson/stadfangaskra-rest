package stadfangaskra

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/StefanKjartansson/stadfangaskra"
	"github.com/golang/groupcache/lru"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/llimllib/loglevel"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidId = errors.New("Invalid id")
)

type HandleFuncErrorStatus func(http.ResponseWriter, *http.Request) (error, int)

type LocationService struct {
	cache          *lru.Cache
	prefix         string
	streetNames    []string
	postCodes      []int
	locations      []stadfangaskra.Location
	indexTable     map[int]*stadfangaskra.Location
	defaultHeaders map[string]string
}

func parseFilter(req *http.Request) (*stadfangaskra.Filter, error) {

	req.ParseForm()
	decoder := schema.NewDecoder()
	f := new(stadfangaskra.Filter)
	err := decoder.Decode(f, req.Form)
	return f, err

}

func NewLocationService(prefix string, locs []stadfangaskra.Location) *LocationService {

	l := LocationService{
		cache:      lru.New(100),
		prefix:     strings.TrimRight(prefix, "/"),
		indexTable: make(map[int]*stadfangaskra.Location),
		locations:  locs,
		defaultHeaders: map[string]string{
			"Content-Type":                 "application/json; charset=utf-8",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		},
	}

	start := time.Now()

	pc := make(map[int]struct{})
	sn := make(map[string]struct{})

	for idx, i := range l.locations {
		l.indexTable[i.ID] = &l.locations[idx]
		pc[i.Postcode] = struct{}{}
		sn[i.Street] = struct{}{}
	}

	for p, _ := range pc {
		l.postCodes = append(l.postCodes, p)
	}

	for k, _ := range sn {
		l.streetNames = append(l.streetNames, k)
	}

	sort.Ints(l.postCodes)
	sort.Strings(l.streetNames)

	log.Infof("Initialize took: %f.ms", time.Now().Sub(start).Seconds()*1000)

	return &l

}

// Error catching middleware
func (l *LocationService) wrapHttpHandler(f HandleFuncErrorStatus) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		for key, value := range l.defaultHeaders {
			w.Header().Set(key, value)
		}
		//start := time.Now()
		err, status := f(w, r)
		if err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), status)
		} else {
			/*
				log.Infof("%s %s %s, time: %f.ms", r.RemoteAddr,
					r.Method, r.URL.Query(),
					time.Now().Sub(start).Seconds()*1000)
			*/
		}

	}
}

func (l *LocationService) listing(w http.ResponseWriter, req *http.Request) (error, int) {

	f, err := parseFilter(req)
	if err != nil {
		return err, http.StatusBadRequest
	}

	key := f.Hash()
	cached, ok := l.cache.Get(key)

	if !ok {
		hasWritten := false
		var b bytes.Buffer
		b.Write([]byte("["))

		for _, l := range l.locations {
			if f.Match(&l) {
				if hasWritten {
					b.Write([]byte(","))
				}
				b.Write(l.JSONCache)
				hasWritten = true
			}
		}
		b.Write([]byte("]"))
		cached = b.Bytes()
		l.cache.Add(key, cached)
	}

	w.Write(cached.([]byte))

	return nil, http.StatusOK
}

func (l *LocationService) detail(w http.ResponseWriter, req *http.Request) (error, int) {

	vars := mux.Vars(req)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return ErrInvalidId, http.StatusBadRequest
	}

	enc := json.NewEncoder(w)
	enc.Encode(l.indexTable[id])

	return nil, http.StatusOK

}

func (l *LocationService) search(w http.ResponseWriter, req *http.Request) (error, int) {

	//enc := json.NewEncoder(w)
	//enc.Encode(l.indexTable[id])

	return nil, http.StatusOK

}

func (l *LocationService) GetRouter() *mux.Router {

	router := mux.NewRouter()
	s := router.PathPrefix(l.prefix).Subrouter()
	s.HandleFunc("/", l.wrapHttpHandler(l.listing)).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}/", l.wrapHttpHandler(l.detail)).Methods("GET")
	return router

}
