FROM golang:1.23-alpine as deps

WORKDIR /app

RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@v1.21.2
RUN CGO_ENABLED=0 go install github.com/cosmtrek/air@v1.49.0
RUN CGO_ENABLED=0 go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

FROM golang:1.23-alpine as final

WORKDIR /app

COPY --from=deps /go/bin /go/bin

COPY go.mod go.sum ./
COPY air.toml ./

RUN go mod tidy

COPY . .

EXPOSE 2345

ENTRYPOINT ["air", "-c", "air.toml"]