FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /usr/app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /usr/dist

RUN cp /build/main .

ENV TOKEN=TOKEN
ENV FIREBASE_CONFIG=FIREBASE_CONFIG
ENV FIREBASE_PROJECTID=FIREBASE_PROJECT_ID
ENV DATADOG_API_KEY=DATADOG_API_KEY

EXPOSE 4000

CMD ["/usr/main"]