FROM alpine:3.16@sha256:65a2763f593ae85fab3b5406dc9e80f744ec5b449f269b699b5efd37a07ad32e as alpine

ARG TARGETARCH

RUN set -eux; \
    apk add -U --no-cache ca-certificates


FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY dist/chekr_${TARGETOS}_${TARGETARCH}*/chekr /usr/local/bin/chekr

ENTRYPOINT ["/usr/local/bin/chekr"]
