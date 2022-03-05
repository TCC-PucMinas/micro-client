# build stage
FROM golang:alpine as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN apk --no-cache add ca-certificates

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o micro-client

# final stage
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/micro-client /app/
ENTRYPOINT ["/app/micro-client"]
