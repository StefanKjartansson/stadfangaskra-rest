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
    "time"
    iconv "github.com/djimenez/iconv-go"
    isnet "./isnet"
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


var placename_file = flag.String("file", "Stadfangaskra_20130326.dsv", "csv input file")
var locations []Location
const shortForm = "02.01.2006"


func RecordToLocation(record []string) (loc Location) {

    length := len(record)

    loc.Hnitnum, _ = strconv.ParseInt(record[0], 0, 64)
    loc.Svfnr = record[1]
    loc.Byggd = record[2]
    loc.Landnr, _ = strconv.ParseInt(record[3], 0, 64)
    loc.Heinum, _ = strconv.ParseInt(record[4], 0, 64)
    loc.Fasteignaheiti = record[5]
    loc.Matsnr = record[6]
    loc.Postnr, _ = strconv.ParseInt(record[7], 0, 64)
    loc.Heiti_Nf = record[8]
    loc.Heiti_Tgf = record[9]
    loc.Husnr, _ = strconv.ParseInt(record[10], 0, 64)
    loc.Bokst = record[11]
    loc.Vidsk = record[12]
    loc.Serheiti = record[13]

    loc.Dags_Inn, _ = time.Parse(shortForm, record[14])
    loc.Dags_Leidr, _ = time.Parse(shortForm, record[15])

    x, _ := strconv.ParseFloat(record[length - 2], 64)
    y, _ := strconv.ParseFloat(record[length - 1], 64)
    loc.X, loc.Y = isnet.Isnet93ToWgs84(x, y)

    return
}

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
        locations = append(locations, RecordToLocation(record))
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
