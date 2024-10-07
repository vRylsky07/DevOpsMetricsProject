FROM golang:latest AS builder
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
COPY cmd/agent/. .
COPY ../../. .
RUN CGO_ENABLED=0 GOOS=linux go build -o agentrun .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/agentrun .
ENTRYPOINT ["./agentrun"]