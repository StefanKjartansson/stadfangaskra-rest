package main

import (
    "testing"
    "net/url"
    "fmt"
)

func TestParseQueryParams(t *testing.T) {

    out1 := []int64{101,200}
    out2 := "Lauga*"
    out3 := []int64{1,10}

    v := url.Values{}
    v.Set("name_startswith", "Lauga")
    v.Add("postcode", "101")
    v.Add("postcode", "200")
    v.Add("number", "1")
    v.Add("number", "10")

    postcodes, numbers, query, err := ParseQueryParams(v)

    if err != nil {
        t.Error(err)
    }

    if fmt.Sprintf("%v", postcodes) != fmt.Sprintf("%v", out1) ||
       fmt.Sprintf("%v", numbers) != fmt.Sprintf("%v", out3) ||
       query != out2 {
		t.Errorf("ParseQueryParams(%v) was (%v,%v), expected (%v,%v)",
            v, postcodes, query, out1, out2)
    }

}

func TestSearchComparison(t *testing.T) {

    loc := Location{
        Postnr:101,
    }

    if !loc.ContainsPostcode([]int64{101}) {
        t.Errorf("Location should contain a found postcode: %v", loc)
    }

    loc.Heiti_Nf = "Laugavegur"
    loc.Heiti_Tgf = "Laugavegi"

    if !loc.NameMatches("Laug*") {
        t.Errorf("Location should name should match: %v", loc)
    }

    if !loc.NameMatches("*vegur") {
        t.Errorf("Location should name should match: %v", loc)
    }

    if !loc.NameMatches("Laugavegur") {
        t.Errorf("Location should name should match: %v", loc)
    }

    if !loc.NameMatches("Laugavegi") {
        t.Errorf("Location should name should match: %v", loc)
    }

}
