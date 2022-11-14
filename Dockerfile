FROM alpine

COPY certs /etc/ssl/certs/
COPY bin/linux-amd64/koroche /
COPY charts/config.yaml /etc/config/

CMD ["/koroche"]
