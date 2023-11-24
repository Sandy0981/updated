FROM golang:1.21.4-alpine3.18 AS builder

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o server cmd/job-portal-api/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app

COPY --from=builder /app/server .

COPY --from=builder /app/private.pem .

COPY --from=builder /app/public.pem .

CMD [ "./server" ]