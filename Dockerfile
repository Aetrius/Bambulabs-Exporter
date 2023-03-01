FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY *.go ./
COPY .env .env

RUN go mod download
RUN go get github.com/sirupsen/logrus
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go get github.com/prometheus/common/version

RUN go mod vendor
RUN go mod tidy
RUN go mod vendor

RUN go build -o /bambulabs-aetrius-exporter

EXPOSE 9101

CMD [ "/bambulabs-aetrius-exporter" ]
