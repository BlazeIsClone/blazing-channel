FROM golang:1.26.2-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api cmd/api/main.go

FROM alpine:3.20.1 AS prod

WORKDIR /app

COPY --from=build /app/api /app/api

EXPOSE ${PORT}

CMD ["./api"]