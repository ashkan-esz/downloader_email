FROM golang:1.21-alpine as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /myapp cmd/main.go

FROM scratch
COPY --from=builder /myapp /myapp

EXPOSE 8888
ENTRYPOINT [ "/myapp" ]
