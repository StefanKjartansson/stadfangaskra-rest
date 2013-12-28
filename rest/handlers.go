package rest

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
	"strconv"
	"strings"
)

var (
	ErrInvalidId = errors.New("Invalid id")
)

type HandleFuncErrorStatus func(http.ResponseWriter, *http.Request) (error, int)

type LocationService struct {
	cache          *lru.Cache
	Store          *stadfangaskra.Store
	prefix         string
	defaultHeaders map[string]string
}

func parseFilter(req *http.Request) (*stadfangaskra.Filter, error) {

	req.ParseForm()
	decoder := schema.NewDecoder()
	f := new(stadfangaskra.Filter)
	err := decoder.Decode(f, req.Form)
	return f, err

}

func NewLocationService(prefix string) *LocationService {

	return &LocationService{
		cache:  lru.New(100),
		prefix: strings.TrimRight(prefix, "/"),
		Store:  stadfangaskra.DefaultStore,
		defaultHeaders: map[string]string{
			"Content-Type":                 "application/json; charset=utf-8",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Origin, X-Requested-With, Content-Type, Accept",
		},
	}

}

// Error catching middleware
func (l *LocationService) wrapHttpHandler(f HandleFuncErrorStatus) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		for key, value := range l.defaultHeaders {
			w.Header().Set(key, value)
		}
		err, status := f(w, r)
		if err != nil {
			log.Errorln(err.Error())
			http.Error(w, err.Error(), status)
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

		for _, l := range l.Store.Locations {
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
	enc.Encode(l.Store.GetById(id))

	return nil, http.StatusOK

}

func (l *LocationService) search(w http.ResponseWriter, req *http.Request) (error, int) {

	val, ok := req.URL.Query()["q"]
	if !ok || len(val) > 1 {
		return nil, http.StatusBadRequest
	}

	enc := json.NewEncoder(w)

	loc, err := l.Store.FindByString(val[0])

	if err != nil {
		return err, http.StatusBadRequest
	}

	enc.Encode(loc)

	return nil, http.StatusOK

}

func (l *LocationService) GetRouter() *mux.Router {

	router := mux.NewRouter()
	s := router.PathPrefix(l.prefix).Subrouter()
	s.HandleFunc("/", l.wrapHttpHandler(l.listing)).Methods("GET")
	s.HandleFunc("/search", l.wrapHttpHandler(l.search)).Methods("GET")
	s.HandleFunc("/{id:[0-9]+}/", l.wrapHttpHandler(l.detail)).Methods("GET")
	return router

}
