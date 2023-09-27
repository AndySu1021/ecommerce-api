FROM golang:1.18 as build

WORKDIR /app

COPY . .

RUN go build -o server main.go

FROM gcr.io/distroless/base

WORKDIR /app

COPY --from=build /app/server /app/server
COPY --from=build /app/db/migrations /app/migrations
COPY --from=build /app/internal/email/templates /app/email/templates
COPY --from=build /app/keys /app/keys

EXPOSE 8082

CMD ["./server", "server"]