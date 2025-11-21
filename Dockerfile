FROM alpine:3.14

RUN apk add --no-cache rust=1.91.1-r0

WORKDIR /app

COPY ./backend/target/release/backend .
