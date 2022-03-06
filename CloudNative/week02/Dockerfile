FROM golang:1.17.8-alpine3.15 AS build

COPY . /src

WORKDIR /src

RUN go mod tidy \
    && CGO_ENABLED=0 go build -o httpserver main.go \
    && chmod +x httpserver

FROM alpine:3.15

COPY --from=build /src/httpserver /httpserver

ENTRYPOINT ["/httpserver"]