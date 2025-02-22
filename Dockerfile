FROM golang:1.23.6

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go install github.com/air-verse/air@latest

CMD ["air"]