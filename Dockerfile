FROM golang:1.4.2-wheezy

ENV APP github.com/jllopis/try5
RUN go get github.com/tools/godep
ADD . /go/src/${APP}
WORKDIR /go/src/${APP}
RUN godep go install ${APP}/cmd/try5d
ADD ./cmd/try5d/certs/*.pem /etc/try5/certs/

EXPOSE 8000
ENTRYPOINT /go/bin/try5d
