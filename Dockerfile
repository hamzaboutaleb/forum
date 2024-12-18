# Build Stage
FROM golang:1.22-alpine

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN CGO_ENABLED=1 go build -o forum .

LABEL Name=forum Version=0.0.1

LABEL Description="This is the Forum app"

EXPOSE 8000

CMD ["./forum"]
