package main

import (
	"testing"
)

func TestDataBase(t *testing.T) {

    const placenameFile = "Stadfangaskra_20130326.dsv"
	ImportDatabase(placenameFile)

	if len(Locations) < 1 {
        t.Errorf("Locations should be larger than 0, is %d.\n",
            len(Locations))
	}
}
