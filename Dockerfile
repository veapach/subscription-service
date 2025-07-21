FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
COPY ./docs ./docs

RUN go build -o main ./cmd/main.go

CMD ["./main"]
