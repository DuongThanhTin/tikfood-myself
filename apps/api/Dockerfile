FROM golang:1.23-alpine AS build

WORKDIR /src

COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go test ./...
RUN go build -o /out/tikfood-api ./cmd/server

FROM alpine:3.20

RUN adduser -D -H appuser
USER appuser

WORKDIR /app
COPY --from=build /out/tikfood-api /app/tikfood-api

EXPOSE 18081

CMD ["/app/tikfood-api"]
