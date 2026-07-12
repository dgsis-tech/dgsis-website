FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build \
	-o dgsis-website \
	./cmd/web


FROM alpine:latest

WORKDIR /app

RUN adduser -D appuser

COPY --from=builder /app/dgsis-website .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

USER appuser

EXPOSE 8080

CMD ["./dgsis-website"]