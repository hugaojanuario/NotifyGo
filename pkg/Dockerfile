FROM golang:1.26.1

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -0 app ./cmd/api

EXPOSE 8080

CMD ["./app"]