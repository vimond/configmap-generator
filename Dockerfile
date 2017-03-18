FROM golang:1.8-alpine
#FROM alpine:3.5

MAINTAINER Dag Viggo Lokoeen <dagviggo@vimond.com>
COPY . /go/src/app
RUN apk add --no-cache --virtual build-dependencies make git
RUN cd /go/src/app && make
RUN apk del build-dependencies
CMD ["help"]
WORKDIR /go/src/app
ENTRYPOINT ["go", "run", "/go/src/app/cmd/cmapgen/main.go"]
