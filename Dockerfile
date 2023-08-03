FROM golang:alpine

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

WORKDIR /app/cmd/gophermart
RUN go build -o /build
EXPOSE 8080

WORKDIR /
CMD ./build -r http://accrual:8090