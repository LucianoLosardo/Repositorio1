FROM golang:latest AS builder
RUN apt-get update
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM scratch
WORKDIR /
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
ENTRYPOINT [ "./main" ]

