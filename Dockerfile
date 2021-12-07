FROM golang:alpine AS builder
WORKDIR /usr/src/app
RUN apk add make
COPY ./go.sum ./go.mod ./Makefile ./
RUN go mod download
COPY . .
RUN make build-alpine

FROM alpine:latest AS runtime
LABEL "org.opencontainers.image.source"="https://github.com/Suggester/suggester-gateway"
WORKDIR /
COPY --from=builder /usr/src/app/suggester-gateway .
CMD ["./suggester-gateway"]
