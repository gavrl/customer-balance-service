FROM golang:1.17-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql, golang-migrate
#RUN curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -
#RUN echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list
RUN apt-get update
RUN apt-get -y install postgresql-client
#RUN apt-get install -y migrate

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

## Substitute BIN for your bin directory.
## Substitute VERSION for the current released version.
## Substitute BINARY_NAME for "buf", "protoc-gen-buf-breaking", or "protoc-gen-buf-lint".
#RUN BIN="/usr/local/bin" && \
#VERSION="0.56.0" && \
#BINARY_NAME="buf" && \
#  curl -sSL \
#    "https://github.com/bufbuild/buf/releases/download/v${VERSION}/${BINARY_NAME}-$(uname -s)-$(uname -m)" \
#    -o "${BIN}/${BINARY_NAME}" && \
#  chmod +x "${BIN}/${BINARY_NAME}"
#
## install Protocol Buffer Compiler
#RUN apt-get -y install protobuf-compiler
#
## install the protocol compiler plugins for Go
#RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
#RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
#
## install gRPC-Gateway
#RUN go install  \
#    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
#    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
#    google.golang.org/protobuf/cmd/protoc-gen-go \
#    google.golang.org/grpc/cmd/protoc-gen-go-grpc

# update your PATH so that the protoc compiler can find the plugins:
RUN export PATH="$PATH:$(go env GOPATH)/bin"

# make wait-for-postgres.sh executable
RUN chmod +x cmd/app/wait-for-postgres.sh

# build go app
#RUN go mod download
RUN go build -o balance-api ./cmd/app/main.go

## generate-proto-go
#RUN buf mod update
#RUN buf build
#RUN buf generate -v --path=api/v1

EXPOSE 8080

CMD ["./balance-api"]