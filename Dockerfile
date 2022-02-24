FROM golang:1.17-alpine as build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/cy-celcat/main.go

# RUN chown -R cycelcat:cycelcat /usr/src/app

EXPOSE 8080

FROM alpine

WORKDIR /cycelcat

RUN apk --no-cache add chromium && \
    adduser -u 1000 -h /cycelcat -D cycelcat

COPY --from=build --chown=cycelcat /usr/local/bin/app /usr/local/bin/app


ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/


USER cycelcat
ENTRYPOINT ["app", "-html=1", "-png=1", "-web=1", "-loop=1"]
