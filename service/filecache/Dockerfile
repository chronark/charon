FROM golang:1.13 as service_builder

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o filecache-service .


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
COPY --from=service_builder /build/filecache-service .
COPY --from=util_builder /waitfor/waitfor .

CMD ["./filecache-service"]

