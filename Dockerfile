FROM golang:1.16-alpine as build

COPY ./ /go/src/github.com/mect/expino-api
WORKDIR /go/src/github.com/mect/expino-api

RUN go build ./cmd/expino-api

FROM alpine:3.13

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/mect/expino-api/expino-api /usr/local/bin/expino-api

RUN mkdir expino-static

CMD ["/usr/local/bin/expino-api", "serve"]