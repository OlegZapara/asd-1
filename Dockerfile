FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o /app/main cmd/sort/main.go

FROM alpine:3.12.0

COPY --from=builder /app/files/in/ /files/in/

COPY --from=builder /app/files/out/ /files/out/

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]

