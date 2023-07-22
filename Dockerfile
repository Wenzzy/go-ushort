# Start from golang base image
FROM golang:1.20-alpine as builder

# Install git.
RUN apk update && apk add --no-cache git

# Working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy everythings
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o ./dist/app-binary ./app/cmd/app/main.go

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app

ADD https://github.com/pressly/goose/releases/download/v3.7.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose
COPY --from=builder /app/dist/app-binary ./
COPY ./app/common/database/migrations ./migrations
COPY ./scripts/prod.sh ./run.sh


#Command to run the executable
CMD ["./app-binary"]
