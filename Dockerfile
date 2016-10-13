FROM alpine:4.3

COPY ./amp-agent /amp-agent

HEALTHCHECK --interval=10s --timeout=15s --retries=12 CMD /amp-agent healthcheck

CMD[/amp-agent]

