## BUILDING

FROM golang:1.12 as build

RUN go get -u github.com/golang/dep/cmd/dep && \
  go get -u github.com/onsi/ginkgo/ginkgo && \
  go get -u github.com/onsi/gomega/...

WORKDIR /go/src/github.com/crazyfacka/yanbapp-api
COPY . .

RUN dep ensure -v -vendor-only && \
  ginkgo ./... && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-w' -v -o /root/app .

## PACKAGING

FROM alpine:3.8

RUN apk add --no-cache ca-certificates

COPY --from=build /root/app /root/app

WORKDIR /root

CMD ["./app"]