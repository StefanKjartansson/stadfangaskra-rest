package stadfangaskra

import (
	isnet "./isnet"
	"encoding/csv"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io"
	"os"
	"strconv"
	"time"
)

func ImportFromRecord(record []string) (loc Location) {

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

	x, _ := strconv.ParseFloat(record[length-2], floatSize)
	y, _ := strconv.ParseFloat(record[length-1], floatSize)
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

	x, _ := iconv.NewReader(file, "iso-8859-1", "utf-8")
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
		Locations = append(Locations, ImportFromRecord(record))
	}
	return
}
