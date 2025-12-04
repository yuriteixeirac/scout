FROM golang:latest

WORKDIR /app

COPY . .

EXPOSE 8080

ENTRYPOINT [ "go", "run", "." ]