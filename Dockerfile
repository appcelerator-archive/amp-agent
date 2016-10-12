FROM alpine:3.4

COPY amp-agent /amp-agent

CMD ["/amp-agent"]