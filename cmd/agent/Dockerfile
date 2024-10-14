FROM golang:latest AS builder
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY cmd/agent/. .
COPY ../../. .
RUN CGO_ENABLED=0 GOOS=linux go build -o DompAgent .

FROM alpine:latest
RUN apk --no-cache add curl
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/DompAgent .
ENTRYPOINT ["./DompAgent"]
