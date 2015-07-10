FROM golang

# Fetch dependencies
RUN go get github.com/tools/godep

# Add project directory to Docker image.
ADD . /go/src/github.com/tbruyelle/go-bootstrap

ENV USER tbruyelle
ENV HTTP_ADDR 8888
ENV HTTP_DRAIN_INTERVAL 1s
ENV COOKIE_SECRET kvLePSAQkjetRLh3

# Replace this with actual PostgreSQL DSN.
ENV DSN postgres://tbruyelle@localhost:5432/go-bootstrap?sslmode=disable

WORKDIR /go/src/github.com/tbruyelle/go-bootstrap

RUN godep go build

EXPOSE 8888
CMD ./go-bootstrap