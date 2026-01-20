FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY notes/go.mod notes/go.sum ./notes/

COPY pkg ./pkg

WORKDIR /app/notes

RUN go mod download

COPY notes/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/notes/main .

EXPOSE 8103

CMD ["./main"]