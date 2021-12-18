FROM golang:alpine AS builder
ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
RUN go mod tidy

COPY . .
RUN go build -mod=mod -o main .

FROM alpine

WORKDIR /dist

ARG TOKEN
ARG PRODUCTION

ENV TOKEN=$TOKEN
ENV PRODUCTION = $PRODUCTION

COPY --from=builder /build/main /dist

ENTRYPOINT ./main 