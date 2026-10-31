FROM golang:1.27.1-alpine

WORKDIR /app

HEALTHCHECK --interval=60s --timeout=1s --start-period=60s --retries=3 \
    CMD wget -qO- http://localhost:$API_PORT/api/health

RUN apk add --no-cache make
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.4
RUN go install github.com/cespare/reflex@v0.3.2

COPY go.mod go.sum ./
RUN go mod download

COPY Makefile ./
COPY src/ src/

EXPOSE $API_PORT

RUN make build

CMD ["make", "run"]
