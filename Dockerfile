FROM golang:1.26.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fiber-app .

FROM alpine:latest

#for communicate with external api that require cer
RUN apk --no-cache add ca-certificates

RUN apk add --no-cache tzdata

WORKDIR /root/

COPY --from=builder /app/fiber-app .

EXPOSE 3000

CMD ["./fiber-app"]


