FROM golang:latest as builder

WORKDIR /app

ENV GOPROXY=direct

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /url-shortener ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /url-shortener .

EXPOSE 8080

CMD ["/app/url-shortener"]