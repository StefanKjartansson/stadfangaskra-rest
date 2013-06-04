package main

import (
    "encoding/csv"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
    iconv "github.com/djimenez/iconv-go"
)


var placename_file = flag.String("file", "Stadfangaskra_20130326.dsv", "csv input file")
var locations []Location

func ImportDatabase(pfile string) {

    file, err := os.Open(pfile)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    x,_ := iconv.NewReader(file, "iso-8859-1", "utf-8")
    reader := csv.NewReader(x)
    reader.Comma = '|'
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Error:", err)
            return
        }
        loc := Location{}
        loc.ImportFromRecord(record)
        locations = append(locations, loc)
    }
    return
}


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

    for _, element := range locations {

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


func main() {
    flag.Parse()
    log.Println("Starting import")
    ImportDatabase(*placename_file)
    log.Println("Data Imported")
    log.Println("Starting server")
    http.HandleFunc("/locations/", GetLocation)
    http.ListenAndServe(":8080", nil)
}
