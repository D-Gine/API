FROM golang:1.25-alpine

WORKDIR /app

HEALTHCHECK --interval=60s --timeout=1s --start-period=1s --retries=3 \
    CMD wget -qO- http://localhost:$API_PORT/api/health

RUN apk update; apk add make

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
COPY src/ src/

EXPOSE $API_PORT

RUN make build

CMD ["make", "run"]
