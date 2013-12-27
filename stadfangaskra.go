package stadfangaskra

import (
	"log"
	"os"
)

var (
	DefaultStore *Store
)

func init() {

	//flag, os, /usr/share/, fixture.json

	file, err := os.Open("./fixture.json")
	if err != nil {
		log.Fatal(err)
	}

	DefaultStore, err = NewStore(file)

	if err != nil {
		log.Fatal(err)
	}

}
