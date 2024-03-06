FROM golang:latest

WORKDIR /app

COPY server.go .
COPY static/ static/


RUN go build -o server

EXPOSE 8080

CMD ["./server"]