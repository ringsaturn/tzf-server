FROM golang:1.20 AS builder

COPY . /src
WORKDIR /src

RUN go mod download
RUN go build

FROM debian:bookworm

LABEL \
    name = "tzf-server" \
    org.opencontainers.image.description = "Simple longititu&latitude to tzname server" \
    org.opencontainers.image.source = "https://github.com/ringsaturn/tzf-server"

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends ca-certificates netbase tzdata \
    && rm -rf /var/lib/apt/lists/ \
    && apt-get autoremove -y

COPY --from=builder /src/tzf-server /app/

WORKDIR /app

EXPOSE 8080

CMD ["/app/tzf-server"]
