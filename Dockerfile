FROM alpine

EXPOSE 80
EXPOSE 8081
EXPOSE 6560

COPY certs /etc/ssl/certs/
COPY bin/linux-amd64/koroche /
COPY charts/config.yaml /etc/config/

CMD ["/koroche"]
