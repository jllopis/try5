FROM golang:1.4.2-wheezy

ENV APP github.com/jllopis/try5
ADD . /go/src/${APP}
ADD ./cmd/try5d/certs/*.pem /etc/try5/certs/
WORKDIR /go/src/${APP}
RUN go get github.com/tools/godep \
    && godep go install ${APP}/cmd/try5d \
    && mkdir /var/lib/try5

EXPOSE 8000
ENTRYPOINT /go/bin/try5d
