FROM golang:1.4.2-wheezy

ENV APP github.com/jllopis/try5
RUN go get github.com/tools/godep
ADD . /go/src/${APP}
#RUN go get -d -v ./...
WORKDIR /go/src/${APP}
RUN godep go install ${APP}/cmd/try5d
EXPOSE 8000
ENTRYPOINT /go/bin/try5d
