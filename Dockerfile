# Stage 1
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
RUN apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR $GOPATH/src/whiterabbittech/arabian-nights/
COPY . .
# Fetch dependencies using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/arabian-nights
ENTRYPOINT ["/go/bin/arabian-nights"]
# Stage 2
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/arabian-nights /go/bin/arabian-nights
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Run the binary.
ENTRYPOINT ["/go/bin/arabian-nights"]
