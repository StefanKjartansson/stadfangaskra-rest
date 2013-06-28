package main

import (
	"strings"
	"time"
)

const (
	shortForm = "02.01.2006"
	intBase   = 0
	intSize   = 64
	floatSize = 64
)

type Location struct {
	Hnitnum        int     `json:"id"`
	Svfnr          string    `json:"-"`
	Byggd          string    `json:"-"`
	Landnr         int     `json:"land_nr"`
	Heinum         int     `json:"-"`
	Fasteignaheiti string    `json:"display_name"`
	Matsnr         string    `json:"-"`
	Postnr         int     `json:"postcode"`
	Heiti_Nf       string    `json:"name_nominative"`
	Heiti_Tgf      string    `json:"name_genitive"`
	Husnr          int     `json:"house_number"`
	Bokst          string    `json:"house_characters,omitempty"`
	Vidsk          string    `json:"suffix,omitempty"`
	Serheiti       string    `json:"special_name,omitempty"`
	Dags_Inn       time.Time `json:"date_added"`
	Dags_Leidr     time.Time `json:"date_update,omitempty"`
	Gagna_Eign     string    `json:"data_owner,omitempty"`
	X              float64   `json:"x"`
	Y              float64   `json:"y"`
}

func (loc Location) ContainsPostcode(list []int) bool {

    if len(list) == 0 {
        return true
    }

	for _, b := range list {
		if b == loc.Postnr {
			return true
		}
	}
	return false
}

func (loc Location) ContainsNumbers(list []int) bool {

    if len(list) == 0 {
        return true
    }

	for _, b := range list {
		if b == loc.Husnr {
			return true
		}
	}
	return false
}

func (loc Location) NameMatches(query string) bool {

	if query == "" {
		return true
	}

	if strings.HasSuffix(query, "*") {
		v := query[0:strings.Index(query, "*")]
		if strings.HasPrefix(loc.Heiti_Nf, v) ||
			strings.HasPrefix(loc.Heiti_Tgf, v) {
			return true
		}
	} else if strings.HasPrefix(query, "*") {
		v := query[strings.Index(query, "*")+1 : len(query)]
		if strings.HasSuffix(loc.Heiti_Nf, v) ||
			strings.HasSuffix(loc.Heiti_Tgf, v) {
			return true
		}
	} else if loc.Heiti_Nf == query ||
		loc.Heiti_Tgf == query {
		return true
	}

	return false
}
