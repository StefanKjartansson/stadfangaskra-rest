stadfangaskra-rest
==================

In-memory REST service for geocoding Icelandic addresses written in golang

### Why?

Needed a very fast service to geocode addresses and wanted to try golang.

### How does it work?

Loads the contents of the Icelandic placenames csv file from [Opin Gögn](http://gogn.island.is/) into memory and wraps it in a warm HTTP blanket.

### Usage

To start the server:

```bash
./stadfangaskra -file=Stadfangaskra_20130326.dsv
```

Example query:
```bash
curl "http://localhost:8080/locations/?postcode=108&postcode=200&name_endswith=vegi&number=12&number=15&number=8"
```