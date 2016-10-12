FROM alpine:3.4

COPY amp-agent /amp-agent

RUN apk -v add curl

CMD ["/amp-agent"]

HEALTHCHECK --interval=10s --timeout=15s --retries=12 CMD curl localhost:3000/api/v1/health