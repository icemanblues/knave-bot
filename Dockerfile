FROM golang:alpine as builder

# Install git # Git is required for fetching the dependencies
# gcc is required for sqlite
RUN apk add --update --no-cache ca-certificates git gcc musl-dev

# Fetch dependencies
RUN mkdir /build 
WORKDIR /build 
COPY go.mod .
COPY go.sum .
RUN go mod download

# build the go binary
COPY main .
RUN CGO_ENABLED=1 GOOS=linux go build --tags "linux" -a -ldflags '-extldflags "-static"' -o knavebot main

# Run the binary in its own scratch container
FROM scratch
COPY --from=builder /build/knavebot /app/
WORKDIR /app
CMD ["./knavebot"]
