FROM alpine:3.16@sha256:7580ece7963bfa863801466c0a488f11c86f85d9988051a9f9c68cb27f6b7872 as alpine

ARG TARGETARCH

RUN set -eux; \
    apk add -U --no-cache ca-certificates


FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY dist/chekr_${TARGETOS}_${TARGETARCH}/chekr /usr/local/bin/chekr

ENTRYPOINT ["/usr/local/bin/chekr"]
