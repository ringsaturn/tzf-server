FROM golang:1.19 AS builder

COPY . /src
WORKDIR /src

RUN make build

FROM debian:stable-slim

LABEL name="tzf-server"
LABEL org.opencontainers.image.description "Simple longititu&latitude to tzname server"
LABEL org.opencontainers.image.source https://github.com/ringsaturn/tzf-server

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/tzf-server /app/
COPY --from=builder /src/info.html /app/

WORKDIR /app

EXPOSE 8080

CMD ["/app/tzf-server"]
