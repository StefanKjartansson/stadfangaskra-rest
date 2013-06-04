package main

import (
	"encoding/csv"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io"
	"os"
)

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
		loc := Location{}
		loc.ImportFromRecord(record)
		locations = append(locations, loc)
	}
	return
}
