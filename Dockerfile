FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

WORKDIR /build/cmd/web
COPY ./cmd/web .
WORKDIR /build
RUN go build -o lidi-server ./cmd/web

WORKDIR /dist
RUN cp /build/lidi-server .
EXPOSE 8080
CMD ["/dist/lidi-server"]