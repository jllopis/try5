FROM golang:1.4.2-wheezy

COPY . /go/src/app
RUN go get -d -v ./...
WORKDIR /go/src/app/server
RUN go install -v
CMD ["app"]
