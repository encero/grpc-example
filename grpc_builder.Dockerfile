FROM golang:1.15 AS plugins
RUN go get github.com/golang/protobuf/protoc-gen-go

FROM alpine:3.12

ARG PROTO_VERSION="3.13.0"

RUN apk --no-cache add protobuf
COPY --from=plugins /go/bin/protoc-gen-go /usr/bin/
