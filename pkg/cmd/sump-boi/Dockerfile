FROM ubuntu:20.04 AS bullshit
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update
RUN apt-get install -y tzdata ca-certificates


FROM scratch
COPY ./bin/sump-boi ./sump-boi
COPY --from=bullshit /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=bullshit /usr/share/zoneinfo /usr/share/zoneinfo
VOLUME ["/data", "/tmp/gpio"]
CMD ["./sump-boi"]