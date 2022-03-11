FROM alpine

RUN apk add linux-headers musl-dev gcc go libpcap-dev ca-certificates git

WORKDIR /mnt

RUN rm -rf /mnt/http_header_capture
ENTRYPOINT ["sh", "build_static.sh"]
