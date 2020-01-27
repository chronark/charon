FROM golang:1.13 as client_builder

WORKDIR /build

RUN go mod init github.com/chronark/charon/client/tiles
RUN go mod download
COPY . .

ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o tiles-client .

FROM alpine:3.11

WORKDIR /
COPY --from=client_builder /build/tiles-client .

CMD ["./tiles-client"]