FROM golang:1.17-alpine

RUN adduser -u 1000 -h /usr/src/app -D cycelcat

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./example.config.json /usr/src/app/config.json

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/cy-celcat/main.go

RUN apk add chromium

EXPOSE 8080

RUN chown -R cycelcat:cycelcat /usr/src/app

ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

USER cycelcat
ENTRYPOINT ["app", "-html=1", "-png=1", "-web=1", "-loop=1"]
