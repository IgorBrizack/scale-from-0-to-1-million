FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o /main

EXPOSE 8020

CMD [ "/main" ]