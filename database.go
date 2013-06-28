package main

import (
    "log"
    isnet "./isnet"
    "encoding/csv"
    iconv "github.com/djimenez/iconv-go"
    "io"
    "os"
    "strconv"
    "time"
)

func ImportFromRecord(record []string) (loc Location, err error) {

    for idx, i := range record {

        switch idx {

            case 0,3,4,7,10:
                val, err := strconv.Atoi(i)
                if err != nil {
                    val = 0
                }
                switch idx {
                    case 0: loc.Hnitnum = val
                    case 3: loc.Landnr = val
                    case 4: loc.Heinum = val
                    case 7: loc.Postnr = val
                    case 10: loc.Husnr = val
                }

            case 14,15:
                val, _ := time.Parse(shortForm, i)
                switch idx {
                    case 14: loc.Dags_Inn = val
                    case 15: loc.Dags_Leidr = val
                }

            case 1: loc.Svfnr = i
            case 2: loc.Byggd = i
            case 5: loc.Fasteignaheiti = i
            case 6: loc.Matsnr = i
            case 8: loc.Heiti_Nf = i
            case 9: loc.Heiti_Tgf = i
            case 11: loc.Bokst = i
            case 12: loc.Vidsk = i
            case 13: loc.Serheiti = i

        }
    }

    x, _ := strconv.ParseFloat(record[22], floatSize)
    y, _ := strconv.ParseFloat(record[23], floatSize)
    loc.X, loc.Y = isnet.Isnet93ToWgs84(x, y)

    return
}

func fileReader(filename string) chan Location {

    buffer := make(chan Location, 256)

    go func() {

        file, err := os.Open(filename)
        if err != nil {
            log.Fatal(err)
        } else {

            x, _ := iconv.NewReader(file, "iso-8859-1", "utf-8")
            reader := csv.NewReader(x)
            reader.Comma = '|'
            _, _ = reader.Read()

            for {
                record, err := reader.Read()
                if err == io.EOF {
                    break
                } else if err != nil {
                    log.Fatal(err)
                }
                t, err := ImportFromRecord(record)
                if err != nil {
                    log.Fatal(err)
                }
                buffer <- t
            }

            close(buffer)
        }
    }()

    return buffer
}

func ImportDatabase(pfile string) {

    readChan := fileReader(pfile)

    closed := false
    for !closed {
        select {
        case ev, ok := <-readChan:
            closed = !ok
            Locations = append(Locations, ev)
        }
    }

    return
}
