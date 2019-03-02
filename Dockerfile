FROM golang:alpine as builder

# Install git # Git is required for fetching the dependencies
RUN apk update && apk add --no-cache git

# Fetch dependencies # Using go get
# not using modules nor a vendoring system yet
WORKDIR $GOPATH/src/github.com/icemanblues/knave-bot
COPY . .
RUN go get -d -v

# build the go binary
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

# Run the binary in its own scratch container
FROM scratch
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
