FROM appcelerator/alpine:20160928

ENV GOPATH /go
ENV PATH $PATH:/go/bin

COPY ./ /go/src/github.com/appcelerator/amp-agent
RUN apk update && \
    apk -v --virtual build-deps add git make go@community musl-dev && \
    apk -v add go@community && \
    go version && \
    cd /go/src/github.com/appcelerator/amp-agent && \
    go get -u github.com/Masterminds/glide/... && \
    glide install && \
    rm -f ./amp-agent && \
    make install && \
    rm /go/bin/glide && \
    apk del binutils-libs binutils gmp isl libgomp libatomic libgcc pkgconf pkgconfig mpfr3 mpc1 libstdc++ gcc go && \
    apk del build-deps && cd / && rm -rf /go/src /go/pkg /var/cache/apk/* /root/.cache /root/.glide


HEALTHCHECK --interval=10s --timeout=15s --retries=12 CMD /go/bin/amp-agent healthcheck

CMD ["/go/bin/amp-agent"]