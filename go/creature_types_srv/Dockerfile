FROM golang:latest AS builder

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -v -o creature_type_srv ./src/

CMD ["/app/creature_type_srv"]
