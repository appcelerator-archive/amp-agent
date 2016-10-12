FROM scratch

COPY amp-agent /amp-agent

#HEALTHCHECK --interval=10s --timeout=15s --retries=12 CMD curl localhost:3000/api/v1/health

CMD ["/amp-agent"]