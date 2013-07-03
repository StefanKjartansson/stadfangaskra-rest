stadfangaskra-rest
==================

[![Build Status](https://secure.travis-ci.org/StefanKjartansson/stadfangaskra-rest.png)](http://travis-ci.org/StefanKjartansson/stadfangaskra-rest)

In-memory REST service for geocoding Icelandic addresses written in golang

### Why?

Needed a very fast service to geocode addresses and wanted to try golang.

It's very fast, request time ranges from a puny 0.03 ms to 8.ms depending on volume returned.

### How does it work?

Loads the contents of the Icelandic placenames csv file from [Opin Gögn](http://gogn.island.is/) into memory and wraps it in a warm HTTP blanket.

### Usage

To start the server:

```bash
./stadfangaskra -file=Stadfangaskra_20130326.dsv
```

Example query:
```bash
curl "http://localhost:3999/locations/?postcode=101"
curl "http://localhost:3999/locations/?postcode=101&name=Selja*"
curl "http://localhost:3999/locations/?postcode=101&name=Seljavegur"
curl "http://localhost:3999/locations/?postcode=101&name=Seljavegur&number=1"
```
