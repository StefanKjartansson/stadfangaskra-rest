package stadfangaskra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/StefanKjartansson/stadfangaskra"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

var serverAddr string
var once sync.Once
var client *http.Client

func startServer() {
	ls := NewLocationService("/locations/", stadfangaskra.Locations)
	http.Handle("/", ls.GetRouter())
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
	client = http.DefaultClient
}

func makeRequest(t *testing.T, method, url string, v interface{}) *http.Response {

	var r *http.Response

	url = fmt.Sprintf("http://%s/locations%s", serverAddr, url)

	t.Logf("[%s]: %s", method, url)

	switch method {
	case "POST", "PUT":
		buf, err := json.Marshal(v)
		if err != nil {
			t.Errorf("Unable to serialize %v to json", v)
		}

		t.Logf("[%s] JSON: %s", method, string(buf))
		req, err := http.NewRequest(method, url, bytes.NewReader(buf))
		req.Header.Add("Content-Type", "application/json")
		if err != nil {
			t.Fatalf("[%s] %s, error: %v", method, url, err)
		}
		r, err = client.Do(req)
		if err != nil {
			t.Fatalf("Error when posting to %s, error: %v", url, err)
		}
	default:
		r, err := client.Get(url)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}

		if r.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(r.Body)
			t.Fatalf("Wrong status code: %d, body:%s\n", r.StatusCode, string(body))
		}

		dec := json.NewDecoder(r.Body)
		defer r.Body.Close()
		err = dec.Decode(v)
		if err != nil {
			t.Fatalf("Error: %v\n", err)
		}
	}

	return r

}

func getJSON(t *testing.T, url string, v interface{}) *http.Response {
	return makeRequest(t, "GET", url, v)
}

func TestFilterURL(t *testing.T) {

	t.Log("Testing filter")

	once.Do(startServer)

	url := "/?street=Laugavegur&number=2"
	results := []stadfangaskra.Location{}
	getJSON(t, url, &results)

	url = "/?street=Laugavegur&number=22"
	results = []stadfangaskra.Location{}
	getJSON(t, url, &results)

}

func TestSearch(t *testing.T) {

	t.Log("Testing search")

	once.Do(startServer)

	url := "/search?q=Vatnsstígur%203b,%20101%20Reykjavík"

	results := []stadfangaskra.Location{}
	getJSON(t, url, &results)

	if results[0].Street != "Vatnsstígur" {
		t.FailNow()
	}

}
