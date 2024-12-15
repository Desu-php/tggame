FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

# Copy the rest of the application code
COPY . .

RUN go build -o main ./cmd/main.go

ENV PORT=8080

EXPOSE ${PORT}

# Start the application
CMD ["./main"]
