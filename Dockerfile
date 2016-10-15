FROM scratch
COPY amp-agent /amp-agent
ENTRYPOINT ["/amp-agent"]
CMD []
