FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build API server from cmd/api-server/main.go
RUN go build -o analytics-app ./cmd/api-server/main.go

EXPOSE 8080

CMD ["./analytics-app"]
