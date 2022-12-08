FROM golang:1.19-alpine as base

RUN mkdir -p /home/app/go
RUN apk add just

WORKDIR /home/app/go

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN just build

FROM alpine as release
COPY --from=base /home/app/go/bin/www /

CMD ["/www"]

