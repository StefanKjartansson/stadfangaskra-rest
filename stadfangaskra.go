package stadfangaskra

import (
	"flag"
	"log"
	"net/http"
)

var placename_file = flag.String("file", "Stadfangaskra_20130326.dsv", "csv input file")
var Locations []Location

func main() {
	flag.Parse()
	log.Println("Starting import")
	ImportDatabase(*placename_file)
	log.Println("Data Imported")
	log.Println("Starting server")
	http.HandleFunc("/locations/", GetLocation)
	http.ListenAndServe(":8080", nil)
	log.Println("Bye")
}
