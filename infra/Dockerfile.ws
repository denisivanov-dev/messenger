FROM golang:1.24

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY .. .

CMD ["/go/bin/air", "-c", ".air.docker.toml"]