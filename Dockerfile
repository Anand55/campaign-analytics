FROM golang:1.24

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o analytics-app main.go

EXPOSE 8080
CMD ["./analytics-app"]
