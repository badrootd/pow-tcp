FROM golang:1.19-alpine AS client

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN make build

CMD ["./bin/client"]

FROM golang:1.19-alpine AS server

RUN apk add --no-cache make

WORKDIR /app

COPY . .

RUN make build

CMD ["./bin/server"]
