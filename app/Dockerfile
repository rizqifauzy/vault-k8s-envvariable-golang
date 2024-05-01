FROM golang:1.22-alpine

WORKDIR /app

RUN apk update && \
    apk add mysql-client

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY app.go ./
RUN go build -o app

EXPOSE 8080

CMD ["./app"]