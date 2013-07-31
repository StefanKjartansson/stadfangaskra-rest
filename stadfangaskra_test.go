package stadfangaskra

import (
	"testing"
)

func TestSearchComparison(t *testing.T) {

	loc := Location{
		Postnr: 101,
	}

	loc.Heiti_Nf = "Laugavegur"
	loc.Heiti_Tgf = "Laugavegi"

	if !loc.NameMatches("Laug*") {
		t.Errorf("Location should name should match: %v", loc)
	}

	if !loc.NameMatches("*vegur") {
		t.Errorf("Location should name should match: %v", loc)
	}

	if !loc.NameMatches("Laugavegur") {
		t.Errorf("Location should name should match: %v", loc)
	}

	if !loc.NameMatches("Laugavegi") {
		t.Errorf("Location should name should match: %v", loc)
	}

}
