FROM golang:1.20 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

COPY --from=builder /app/main /server

CMD ["/server"]
