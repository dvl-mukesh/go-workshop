FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM alpine:latest AS production
COPY --from=builder /app .

CMD [ "./main" ]


