stadfangaskra-rest
==================

[![Build Status](https://secure.travis-ci.org/StefanKjartansson/stadfangaskra-rest.png)](http://travis-ci.org/StefanKjartansson/stadfangaskra-rest)

REST service for geocoding Icelandic addresses written in golang. Loads a `json` fixture file into memory containing the Icelandic placename information sourced from the csv file from [Opin Gögn](http://gogn.island.is/).

## Why?

Needed a very fast service to geocode addresses and wanted to try golang.

It's very fast, request time ranges from a puny 0.01 ms to 15.ms depending on volume of data returned.

### Benchmarks

All benchmark ran on a 2012 Macbook Pro, 20 concurrent connections, 4 threads and 30 seconds

#### Large set

All addresses in postcode 101.

```bash
wrk -c 20 -d 30 -t 4 "http://localhost:3999/locations/?postcode=101"
Running 30s test @ http://localhost:3999/locations/?postcode=101
  4 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     4.20ms    5.00ms  44.31ms   97.96%
    Req/Sec     1.44k   269.75     2.18k    77.19%
  164929 requests in 30.00s, 110.00GB read
Requests/sec:   5497.56
Transfer/sec:      3.67GB
```

#### Single address

```bash
wrk -c 20 -d 30 -t 4 "http://localhost:3999/locations/10083841/"
Running 30s test @ http://localhost:3999/locations/10083841/
  4 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.95ms    5.63ms  38.91ms   97.59%
    Req/Sec     4.67k     0.88k    5.67k    92.91%
  538249 requests in 30.00s, 220.73MB read
Requests/sec:  17941.44
Transfer/sec:      7.36MB
```

#### Realistic query

```bash
wrk -c 20 -d 30 -t 4 "http://localhost:3999/locations/?postcode=101&street=Laugavegur&number=1&number=22"
Running 30s test @ http://localhost:3999/locations/?postcode=101&street=Laugavegur&number=1&number=22
  4 threads and 20 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     2.86ms    7.62ms  42.06ms   96.24%
    Req/Sec     3.65k   823.27     4.44k    94.04%
  416528 requests in 30.00s, 430.20MB read
Requests/sec:  13884.14
Transfer/sec:     14.34MB
```

## Usage

Start the server:

```bash
./stadfangaskra fixture.json
```

### Example queries:

```bash
curl "http://localhost:3999/locations/?postcode=101"
curl "http://localhost:3999/locations/?postcode=101&street=Selja*"
curl "http://localhost:3999/locations/?postcode=101&street=Seljavegur"
curl "http://localhost:3999/locations/?postcode=101&street=Seljavegur&number=1"
```
