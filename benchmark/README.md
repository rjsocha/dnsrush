## Distributed DNS benchmark

Quick 
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
cat example/example.play | docker exec -i dnsrushd upload

# Deploy worker nodes 
docker exec dnsrushd deploy

# List connected nodes
docker exec dnsrushd list

# Run benchmark
docker exec dnsrushd play 1000

# get result from node
docker exec dnsrushd result <node-id>
```
