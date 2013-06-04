package main

import (
    "strconv"
    "strings"
    "time"
    isnet "./isnet"
)

const (
    shortForm = "02.01.2006"
    intBase = 0
    intSize = 64
    floatSize = 64
)


type Location struct {

    Hnitnum int64 `json:"id"`
    Svfnr string `json:"-"`
    Byggd string `json:"-"`
    Landnr int64 `json:"land_nr"`
    Heinum int64 `json:"-"`
    Fasteignaheiti string `json:"display_name"`
    Matsnr string `json:"-"`
    Postnr int64 `json:"postcode"`
    Heiti_Nf string `json:"name_nominative"`
    Heiti_Tgf string `json:"name_genitive"`
    Husnr int64 `json:"house_number"`
    Bokst string `json:"house_characters,omitempty"`
    Vidsk string `json:"suffix,omitempty"`
    Serheiti string `json:"special_name,omitempty"`
    Dags_Inn time.Time `json:"date_added"`
    Dags_Leidr time.Time `json:"date_update,omitempty"`
    Gagna_Eign string `json:"data_owner,omitempty"`
    X float64 `json:"x"`
    Y float64 `json:"y"`
}


func (loc Location) ImportFromRecord(record []string) {

    length := len(record)

    loc.Hnitnum, _ = strconv.ParseInt(record[0], intBase, intSize)
    loc.Svfnr = record[1]
    loc.Byggd = record[2]
    loc.Landnr, _ = strconv.ParseInt(record[3], intBase, intSize)
    loc.Heinum, _ = strconv.ParseInt(record[4], intBase, intSize)
    loc.Fasteignaheiti = record[5]
    loc.Matsnr = record[6]
    loc.Postnr, _ = strconv.ParseInt(record[7], intBase, intSize)
    loc.Heiti_Nf = record[8]
    loc.Heiti_Tgf = record[9]
    loc.Husnr, _ = strconv.ParseInt(record[10], intBase, intSize)
    loc.Bokst = record[11]
    loc.Vidsk = record[12]
    loc.Serheiti = record[13]

    loc.Dags_Inn, _ = time.Parse(shortForm, record[14])
    loc.Dags_Leidr, _ = time.Parse(shortForm, record[15])

    x, _ := strconv.ParseFloat(record[length - 2], floatSize)
    y, _ := strconv.ParseFloat(record[length - 1], floatSize)
    loc.X, loc.Y = isnet.Isnet93ToWgs84(x, y)

    return
}


func (loc Location) ContainsPostcode (list []int64) bool {
    for _, b := range list {
        if b == loc.Postnr {
            return true
        }
    }
    return false
}

func (loc Location) ContainsNumbers (list []int64) bool {
    for _, b := range list {
        if b == loc.Husnr {
            return true
        }
    }
    return false
}


func (loc Location) NameMatches (query string) bool {

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
        v := query[strings.Index(query, "*") + 1:len(query)]
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
