FROM golang:alpine AS builder
WORKDIR /app
COPY . ./
RUN go build -o app cmd/api-server/*.go
FROM alpine
WORKDIR /app
COPY --from=builder /app/app ./
COPY --from=builder /app/config /app/config
COPY --from=builder /app/Makefile ./
EXPOSE 8000
CMD [ "./app",  "-config=config/docker.yaml" ]