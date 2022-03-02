FROM golang:1.17-alpine
WORKDIR /app
COPY . .
RUN go build -o /bin/mygram_migrate ./cmd/migrate/
RUN go build -o /bin/mygram ./cmd/mygram/
CMD ["/bin/mygram"]
