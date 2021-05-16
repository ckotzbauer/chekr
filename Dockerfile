FROM alpine:3.13

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app
COPY dist/chekr_${TARGETOS}_${TARGETARCH}/chekr .

ENTRYPOINT ["/app/chekr"]
