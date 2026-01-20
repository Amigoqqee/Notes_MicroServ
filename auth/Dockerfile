FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY auth/go.mod auth/go.sum ./auth/

COPY pkg ./pkg

WORKDIR /app/auth

RUN go mod download

COPY auth/ .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/auth/main .

EXPOSE 8101

CMD ["./main"]