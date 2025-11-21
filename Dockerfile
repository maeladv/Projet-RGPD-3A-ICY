FROM golang:1.25

WORKDIR /app

COPY backend/go.mod ./
RUN go mod download

COPY backend/*.go ./
RUN go build -o /srv/www/backend

WORKDIR /srv/www

CMD ["./backend"]
