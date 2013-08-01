package stadfangaskra

import (
	"code.google.com/p/gorilla/mux"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
)

var serverAddr string
var once sync.Once

func startServer() {

	router := new(mux.Router)
	SetupRoutes(router)
	http.Handle("/", router)
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	log.Print("Test Server running on ", serverAddr)

}

func testRequest(t *testing.T, url string, v interface{}) {

	t.Logf("[GET]: %s\n", url)

	r, err := http.Get(url)

	if err != nil {
		t.Errorf("Error: %v\n", err)
	}

	if r.StatusCode != http.StatusOK {
		t.Errorf("Wrong status code: %d\n", r.StatusCode)
	}

	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()

	err = json.Unmarshal(content, v)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}

}

func TestSingleHandler(t *testing.T) {

	once.Do(startServer)

	url := fmt.Sprintf("http://%s/locations/10015125/", serverAddr)
	loc := Location{}

	testRequest(t, url, &loc)

}

func TestAutoCompleteHandler(t *testing.T) {

	once.Do(startServer)

	url := fmt.Sprintf("http://%s/ac/streets/?q=Borg", serverAddr)
	res := []string{}

	testRequest(t, url, &res)

	t.Log(res)

}

func TestSearchHandlers(t *testing.T) {

	once.Do(startServer)

	tests := []struct {
		Params url.Values
		Status int
	}{{
		Params: url.Values{
			"postcode": {"101"},
		},
		Status: http.StatusOK,
	}, {
		Params: url.Values{
			"postcode": {"101"},
			"name":     {"Seljavegur"},
			"number":   {"2"},
		},
		Status: http.StatusOK,
	}, {
		Params: url.Values{
			"name":   {"Seljavegur"},
			"number": {"2"},
		},
		Status: http.StatusOK,
	}, {
		Params: url.Values{
			"name":     {"*vegur"},
			"number":   {"2"},
			"postcode": {"101", "200"},
		},
		Status: http.StatusOK,
	}}

	for _, test := range tests {
		url := fmt.Sprintf("http://%s/locations/?%s", serverAddr, test.Params.Encode())
		locs := []Location{}
		testRequest(t, url, &locs)
	}
}
