FROM arigaio/atlas:latest-community AS atlas

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=atlas /atlas /usr/local/bin/atlas
COPY database/ /schema/
COPY docker/migrate/entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/bin/sh", "/entrypoint.sh"]
