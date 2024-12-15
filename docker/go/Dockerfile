FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN apk add --no-cache make bash build-base && \
    go mod tidy && go mod download && \
    go install github.com/pressly/goose/v3/cmd/goose@latest && \
    echo "export PATH=$PATH:$(go env GOPATH)/bin" >> /etc/profile




# Copy the rest of the application code
COPY . .

RUN go build -o main ./cmd/main.go

ENV PORT=8080

EXPOSE ${PORT}

# Start the application
CMD ["./main"]
