FROM alpine:3.4

MAINTAINER Dag Viggo Lokoeen <dagviggo@vimond.com>

RUN apk add --no-cache ca-certificates
RUN mkdir -p /app
COPY cmapgen /app
CMD ["/app/cmapgen"]