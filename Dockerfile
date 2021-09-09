FROM golang:1.16.7-buster

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x cmd/wait-for-postgres.sh

# build go app
#RUN go mod download
RUN go build -o customer-api ./cmd/main.go

EXPOSE 8080

CMD ["./customer-api"]