FROM golang:1.23-alpine

RUN mkdir /app
RUN mkdir /app/pkg
RUN mkdir /app/cmd
RUN mkdir /app/internal
RUN apk update && apk add --upgrade git openssh

WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY vendor ./vendor

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -mod vendor -a -installsuffix nocgo -o /bin/backend cmd/backend/main.go

FROM alpine:latest  
COPY --from=0 /bin/backend /bin

WORKDIR /
ENTRYPOINT ["/bin/backend"]


