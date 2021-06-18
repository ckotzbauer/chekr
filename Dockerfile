FROM alpine:3.14

ARG TARGETOS
ARG TARGETARCH

RUN addgroup -g 1000 chekr && \
    adduser -u 1000 -G chekr -s /bin/sh -D chekr

COPY dist/chekr_${TARGETOS}_${TARGETARCH}/chekr /usr/local/bin/chekr

ENTRYPOINT ["/usr/local/bin/chekr"]
USER chekr
