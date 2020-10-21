FROM golang:alpine AS builder
ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o main .

FROM alpine

WORKDIR /dist
ARG TOKEN
ARG FIREBASE_CONFIG
ARG FIREBASE_PROJECT_ID
ARG DATADOG_API_KEY
ENV PRODUCTION=TRUE
ENV TOKEN=$TOKEN
ENV FIREBASE_CONFIG=$FIREBASE_CONFIG
ENV FIREBASE_PROJECT_ID=$FIREBASE_PROJECT_ID
ENV DATADOG_API_KEY=$DATADOG_API_KEY
RUN apk update && apk upgrade && apk add --no-cache bash git && apk add --no-cache chromium

# Installs latest Chromium package.
RUN echo @edge http://nl.alpinelinux.org/alpine/edge/community >> /etc/apk/repositories \
    && echo @edge http://nl.alpinelinux.org/alpine/edge/main >> /etc/apk/repositories \
    && apk add --no-cache \
    harfbuzz@edge \
    nss@edge \
    freetype@edge \
    ttf-freefont@edge \
    && rm -rf /var/cache/* \
    && mkdir /var/cache/apk


CMD chromium-browser --headless --disable-gpu --remote-debugging-port=9222 --disable-web-security --safebrowsing-disable-auto-update --disable-sync --disable-default-apps --hide-scrollbars --metrics-recording-only --mute-audio --no-first-run --no-sandbox

COPY --from=builder /build/main /dist
COPY --from=builder /build/resources /dist/resources
COPY --from=builder /build/src/interpreter/libs /dist/src/interpreter/libs

ENTRYPOINT ./main 