FROM golang:1.24.1-alpine

WORKDIR /goapp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/web

ENV PORT=8080

EXPOSE $PORT

CMD [ "sh", "-c", "./main -addr :$PORT"]
