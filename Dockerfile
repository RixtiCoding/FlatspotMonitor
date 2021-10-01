FROM golang:1.17-alpine3.14
RUN mkdir /app1
ADD . /app1
WORKDIR /app1

RUN go mod download

RUN go build -o main .

CMD ["/app1/main"]