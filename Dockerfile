FROM golang:alpine as builder

WORKDIR /usr/local/go/src/

ADD . .

RUN go clean --modcache
RUN go build -mod=readonly -o service cmd/authservice/main.go

FROM alpine:latest

COPY --from=builder /usr/local/go/src/service /

EXPOSE 9000

CMD ["/service"]