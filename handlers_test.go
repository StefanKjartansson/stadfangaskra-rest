package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHandlers(t *testing.T) {

	const path = "/locations/?"

	tests := []struct {
		Desc    string
		Handler func(http.ResponseWriter, *http.Request)
		Path    string
		Params  url.Values
		Status  int
	}{{
		Desc:    "",
		Handler: LocationSearchHandler,
		Path:    path,
		Params: url.Values{
			"postcode": {"101"},
		},
		Status: http.StatusOK,
	}, {
		Desc:    "",
		Handler: LocationSearchHandler,
		Path:    path,
		Params: url.Values{
			"postcode": {"101"},
			"name":     {"Seljavegur"},
			"number":   {"2"},
		},
		Status: http.StatusOK,
	}, {
		Desc:    "",
		Handler: LocationSearchHandler,
		Path:    path,
		Params: url.Values{
			"name":   {"Seljavegur"},
			"number": {"2"},
		},
		Status: http.StatusOK,
	}, {
		Desc:    "Single entity",
		Handler: LocationDetailHandler,
		Path:    "/locations/10015125/",
		Status:  http.StatusOK,
	}}

	for _, test := range tests {
		record := httptest.NewRecorder()
		req, err := http.NewRequest("GET", test.Path+test.Params.Encode(), nil)
		if err != nil {
			t.Fatal(err)
		}
		test.Handler(record, req)
		if got, want := record.Code, test.Status; got != want {
			t.Errorf("%s: response code = %d, want %d", test.Desc, got, want)
		}
	}
}
