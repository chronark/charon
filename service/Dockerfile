FROM golang:1.13 as service_builder
ARG SERVICE

WORKDIR /charon

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
RUN go build -o /${SERVICE} ./service/${SERVICE}


FROM alpine:3.11
ARG SERVICE

WORKDIR /
COPY --from=service_builder /${SERVICE} .
CMD [./${SERVICE}]

