FROM golang:1.21-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go build -o /vulnapi ./cmd/vulnapi

FROM alpine:3.18

RUN apk add --no-cache iputils

WORKDIR /app
COPY --from=build /vulnapi /usr/local/bin/vulnapi

ENV PORT=9000
ENV DB_PATH=/data/vulnapi.db
VOLUME ["/data"]
EXPOSE 9000

CMD ["vulnapi"]
