FROM golang:1.17-alpine as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./example.config.json /usr/src/app/config.json
COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/cy-celcat/main.go

# RUN chown -R cycelcat:cycelcat /usr/src/app

EXPOSE 8080

FROM alpine

RUN apk --no-cache add chromium && \
    mkdir -p /usr/src/app/out && \
    adduser -u 1000 -h /usr/src/app -D cycelcat

COPY --from=build --chown=cycelcat /usr/local/bin/app /usr/local/bin/app


ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

WORKDIR /usr/src/app

USER cycelcat
ENTRYPOINT ["app", "-html=1", "-png=1", "-web=1", "-loop=1"]
