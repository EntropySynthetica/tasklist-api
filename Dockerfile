FROM golang:1.17.0-bullseye

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /tasklist-api

CMD ["/tasklist-api"]

