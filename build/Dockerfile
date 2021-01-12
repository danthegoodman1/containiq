FROM golang:alpine

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o /go/bin/containiq ./cmd/containiq

#Stage 2
From alpine:latest

COPY --from=0 /go/bin/containiq  /go/bin/containiq

#run binary
CMD ["/go/bin/containiq"]