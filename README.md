# Distributed DNS benchmark

License: public domain

Staus: alpha software

For my internal use - so I don't care about code/style/etc.


## Quick Start
```
# Run server
docker run -d --rm -p 80:80 --name dnsrushd rjsocha/dnsrushd

# Configure server hostname
docker exec dnsrushd me <hostname:port>

# Start new benchmark
docker exec dnsrushd new

# Set nameserver
docker exec dnsrushd ns <IP>

# Upload playlist
curl -sfL https://raw.githubusercontent.com/rjsocha/dnsrush/master/example/example.net.play | docker exec -i dnsrushd upload

# Deploy worker nodes 
docker exec dnsrushd deploy

# List connected nodes
docker exec dnsrushd list

# Run benchmark
docker exec dnsrushd play 1000

# get result from node
docker exec dnsrushd result <node-id>
```

# dnsrush

Very simple (qucik&dirty) DNS server benchmark tool in Go.

## Usage of dnsrush tool

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

## Output
```
1591177993659639359 0 0 181523 example.com. 1

start_ts	(query timestamp in nanoseconds)
status 		(0 - OK, 1 - connection error, 2 - query error)
rcode  		(0 - OK)
rtt		(response time in nanoseconds)
query		
query_type
```
