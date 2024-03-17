FROM golang:1.22 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN mkdir -p -m 0700 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts
RUN --mount=type=ssh 
RUN git config --global --add url."ssh://git@github.com".insteadOf "https://github.com"

ENV GOPRIVATE=github.com/Digivate-Labs-Pvt-Ltd/*

RUN --mount=type=ssh CGO_ENABLED=0 GOOS=linux go build -o app cmd/server/main.go

FROM alpine:latest AS production
COPY --from=builder /app .

CMD ["./app"]