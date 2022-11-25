FROM alpine:3.17@sha256:8914eb54f968791faf6a8638949e480fef81e697984fba772b3976835194c6d4 as alpine

ARG TARGETARCH

RUN set -eux; \
    apk add -U --no-cache ca-certificates


FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY dist/chekr_${TARGETOS}_${TARGETARCH}*/chekr /usr/local/bin/chekr

ENTRYPOINT ["/usr/local/bin/chekr"]
