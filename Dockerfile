FROM alpine:1.6

COPY amp-agent /amp-agent

CMD ["/amp-agent"]