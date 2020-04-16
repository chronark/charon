FROM golang:1.14 as util_builder
WORKDIR /

RUN git clone https://github.com/chronark/waitfor.git
WORKDIR /waitfor

RUN go mod download
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o waitfor .

FROM golang:1.14 as service_builder

WORKDIR /charon

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o /tiles-service ./service/tiles


FROM alpine:3.11

WORKDIR /
COPY --from=service_builder /tiles-service .
COPY --from=util_builder /waitfor/waitfor .

CMD ["./tiles-service"]

