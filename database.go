package stadfangaskra

import (
	"encoding/csv"
	"encoding/json"
	"github.com/StefanKjartansson/isnet93"
	iconv "github.com/djimenez/iconv-go"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	Locations        []Location
	IndexTable       = make(map[int]*Location)
	LookupTable      = make(map[[2]int][]*Location)
	DefaultNumbers   = []int{}
	DefaultPostCodes = []int{}
	StreetNames      = []string{}
)

func ImportFromRecord(record []string) (loc Location, err error) {

	for idx, i := range record {

		switch idx {

		case 0, 3, 4, 7, 10:
			val, err := strconv.Atoi(i)
			if err != nil {
				val = 0
			}
			switch idx {
			case 0:
				loc.Hnitnum = val
			case 3:
				loc.Landnr = val
			case 4:
				loc.Heinum = val
			case 7:
				loc.Postnr = val
			case 10:
				loc.Husnr = val
			}

		case 14, 15:
			val, _ := time.Parse(shortForm, i)
			switch idx {
			case 14:
				loc.Dags_Inn = val
			case 15:
				loc.Dags_Leidr = val
			}

		case 1:
			loc.Svfnr = i
		case 2:
			loc.Byggd = i
		case 5:
			loc.Fasteignaheiti = i
		case 6:
			loc.Matsnr = i
		case 8:
			loc.Heiti_Nf = i
		case 9:
			loc.Heiti_Tgf = i
		case 11:
			loc.Bokst = i
		case 12:
			loc.Vidsk = i
		case 13:
			loc.Serheiti = i

		}
	}

	x, _ := strconv.ParseFloat(record[22], floatSize)
	y, _ := strconv.ParseFloat(record[23], floatSize)
	loc.X, loc.Y = isnet93.Isnet93ToWgs84(x, y)

	return
}

func importStream(source io.Reader) error {

	reader := csv.NewReader(source)
	reader.Comma = '|'
	_, _ = reader.Read()

	for {
		r, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		t, err := ImportFromRecord(r)

		if err != nil {
			return err
		}

		Locations = append(Locations, t)
	}

	return nil
}

func postProcess() {

	log.Println("Postprocessing started")

	maxNum := 0
	pnrs := make(map[int]struct{})
	streetnames := make(map[string]struct{})

	for idx, l := range Locations {

		b, err := json.Marshal(l)
		if err != nil {
			log.Fatal(err)
		}
		Locations[idx].JSONCache = b

		key := [2]int{l.Postnr, l.Husnr}
		pnrs[l.Postnr] = struct{}{}
		LookupTable[key] = append(LookupTable[key], &Locations[idx])
		IndexTable[l.Hnitnum] = &Locations[idx]

		streetnames[strings.TrimSpace(l.Heiti_Nf)] = struct{}{}

		if maxNum < l.Husnr {
			maxNum = l.Husnr
		}
	}

	for i := 1; i < maxNum+1; i++ {
		DefaultNumbers = append(DefaultNumbers, i)
	}

	for p, _ := range pnrs {
		DefaultPostCodes = append(DefaultPostCodes, p)
	}

	for k, _ := range streetnames {
		StreetNames = append(StreetNames, k)
	}

	sort.Ints(DefaultPostCodes)
	sort.Strings(StreetNames)

	log.Println("Postprocessing finished")
}

func ImportDatabase(pfile string) {

	log.Println("Starting import")

	defer postProcess()

	file, err := os.Open(pfile)
	if err != nil {
		log.Fatal(err)
	}

	x, _ := iconv.NewReader(file, "iso-8859-1", "utf-8")

	importStream(x)

	log.Println("Import finished")

}
