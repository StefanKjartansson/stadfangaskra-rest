package stadfangaskra

type Location struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Number       int    `json:"house_number,omitempty"`
	NumberChars  string `json:"house_characters,omitempty"`
	SpecificName string `json:"specific_name,omitempty"`
	Street       string `json:"street,omitempty"`
	Postcode     int    `json:"postcode"`
	Municipality string `json:"municipality,omitempty"`
	Coordinates  struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"coordinates"`
	JSONCache []byte `json:"-"`
}
