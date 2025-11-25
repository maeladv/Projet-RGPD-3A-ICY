FROM golang:1.22

WORKDIR /app

# Copier go.mod et go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copier tous les fichiers du backend (handlers, db, models, main.go, etc.)
COPY . ./

# Build l'application
RUN go build -o /srv/www/backend

WORKDIR /srv/www

# Exposer le port
EXPOSE 8080

CMD ["./backend"]