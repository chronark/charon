FROM golang:1.13 as service_builder

WORKDIR /build

RUN go mod init github.com/chronark/charon/service/geocoding
RUN go mod download

COPY . .
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o geocoding-service .


FROM golang:1.13 as util_builder

WORKDIR /

RUN git clone https://github.com/chronark/waitfor.git
WORKDIR /waitfor

RUN go mod download
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o waitfor .

FROM alpine:3.11

WORKDIR /
COPY --from=service_builder /build/geocoding-service .
COPY --from=util_builder /waitfor/waitfor .

CMD ["./geocoding-service"]

