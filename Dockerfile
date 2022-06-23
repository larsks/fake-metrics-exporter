FROM docker.io/golang:1.18 AS builder

COPY . /src
WORKDIR /src
RUN CGO_ENABLED=0 go build

FROM scratch

EXPOSE 9283
COPY --from=builder /src/fake-metrics-exporter /fake-metrics-exporter
ENTRYPOINT ["/fake-metrics-exporter"]
