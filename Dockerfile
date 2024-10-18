FROM golang:1.23.0-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev
#RUN apt-get install bash make git make gcc
#RUN sudo apt install gettext

COPY ["./go.mod", "./go.sum", "./"]

RUN go mod download

COPY . ./
RUN go build -o ./bin/dem3 ./cmd/web/

FROM alpine
COPY --from=builder /usr/local/src/bin/dem3 /
COPY config.yaml /config.yaml
COPY temp /temp
COPY ui /ui

CMD ["/dem3"]