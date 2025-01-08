# STEP 1 - Compile binary

FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO will disable the use of C code. It will help to create a lightweight docker image
# GOOS will set the target operating system to Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o main .


# STEP 2 - Create final image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Set default environment variables
ENV PORT=8080
ENV DATABASE_URL="postgres://postgres:postgres@db:5432/userdb?sslmode=disable"
ENV RANDOM_USER_API="https://randomuser.me/api/"

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]