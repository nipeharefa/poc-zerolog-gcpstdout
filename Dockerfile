# FROM golang:1.17.11-bullseye AS builder
FROM golang:1.18-bullseye as builder

# ENV GO111MODULE=on
ENV GOPRIVATE="dev.azure.com"
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ARG VERSION=development
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
	-ldflags="-w -s" -o build/auth *.go

FROM scratch
# RUN apk update && apk add --no-cache ca-certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/build /app
WORKDIR /app
CMD ["./auth"]