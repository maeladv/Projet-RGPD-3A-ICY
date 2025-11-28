FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /srv/www/backend

WORKDIR /srv/www

CMD ["./backend"]
