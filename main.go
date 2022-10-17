package main

import (
	"bytes"
	"flag"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	addr := flag.String("listen", ":8080", "address to listen to")
	cert := flag.String("cert", "", "server certificate")
	key := flag.String("key", "", "server private key")
	upstream := flag.String("upstream", "https://eu-central.pkg.julialang.org", "upstream server")
	speed := flag.Float64("speed", 1, "speed of body \"processing\" in MB/s")
	flag.Parse()

	upstreamURL, err := url.Parse(*upstream)
	if err != nil {
		log.Fatal(err)
	}

	server := &http.Server{
		Addr: *addr,
		Handler: &httputil.ReverseProxy{
			Director: func(r *http.Request) {
				r.URL.Scheme = upstreamURL.Scheme
				r.URL.Host = upstreamURL.Host
				log.Infof("Request for %s", r.URL.Path)
			},
			ModifyResponse: func(r *http.Response) error {
				log.Infof("Response for %s: %s (%d bytes)", r.Request.URL.Path, r.Status, r.ContentLength)
				if ok, _ := regexp.MatchString("^/(?:artifact|package|registry)/", r.Request.URL.Path); ok && r.Request.Method == http.MethodGet {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						return err
					}
					wait := time.Duration(int64(float64(r.ContentLength)/(1024*1024**speed))) * time.Second
					log.Infof("Body downloaded for %s - waiting %s", r.Request.URL.Path, wait)
					timer := time.NewTimer(wait)
					<-timer.C
					log.Infof("Serving up response for %s", r.Request.URL.Path)
					r.Body = io.NopCloser(bytes.NewReader(body))
				}

				return nil
			},
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}

	log.Infof("Listening on %s", *addr)
	if *cert != "" && *key != "" {
		log.Fatal(server.ListenAndServeTLS(*cert, *key))
	} else {
		log.Fatal(server.ListenAndServe())
	}
}
