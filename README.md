# dnsrush
Very simple (qucik&dirty) DNS server benchmark tool in Go.


# Usage

```
Usage: ./dnsrush
  -c int
        number of queries (default 1)
  -mode string
        playlist mode: random/sequential (default "random")
  -ns string
        nameserver IPv4 address (required)
  -playlist string
        playlist input file
  -q string
        query
  -qt string
        query type
  -t int
        connection timeout (in ms) (default 25)
```


Simple test

```
./dnsrush -ns 127.0.0.1 -q example.com -qt A -c 100
```

Playlist mode
```
./dnsrush -ns 127.0.0.1 -playlist playlist.txt -mode random -c 100

or

./dnsrush -ns 127.0.0.1 -playlist playlist.txt -mode sequential -c 100
```

# Output
```
0 0 1591177993659639359 181523 example.com. 1

status 		(0 - OK, 1 - connection error, 2 - query error)
rcode  		(0 - OK)
start_ts	(query timestamp in nanoseconds)
rtt		(response time in nanoseconds)
query		
query_type
```
