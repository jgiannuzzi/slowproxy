FROM golang:1.19 AS build
WORKDIR /build/slowproxy
COPY go.mod go.sum *.go ./
RUN CGO_ENABLED=0 go build

FROM alpine
COPY --from=build /build/slowproxy/slowproxy /usr/local/bin/
ENTRYPOINT ["slowproxy"]
