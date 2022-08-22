FROM golang:1.18-alpine

# Set the Current Working Directory inside the container
WORKDIR /go/src/money-diff

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go build -o ./out/money-diff .


# Run the binary program produced by `go install`
CMD ["./out/money-diff"]