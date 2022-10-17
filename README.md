### Description
`slowproxy` emulates a reverse proxy that would download the entire response body before scanning it.

### Usage
Build and run with the Go SDK or with Docker.

#### Go
```
go build
./slowproxy
```

#### Docker
```
docker build -t slowproxy .
docker run --rm -ti -p 8080:8080 slowproxy
```

#### Flags
Add `--help` to learn about the various flags.
```
Usage of slowproxy:
  -cert string
    	server certificate
  -key string
    	server private key
  -listen string
    	address to listen to (default ":8080")
  -speed float
    	speed of body "processing" in MB/s (default 1)
  -upstream string
    	upstream server (default "https://eu-central.pkg.julialang.org")
````