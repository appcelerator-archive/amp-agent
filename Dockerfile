FROM appcelerator/alpine:3.4

ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
RUN mkdir -p /go/src/github.com/appcelerator/amp-agent /go/bin
WORKDIR /go/src/github.com/appcelerator/amp-agent

RUN apk update && apk --virtual build-deps add go git curl

COPY ./ ./
RUN go get -u github.com/Masterminds/glide/... && \
    glide install && \
    rm -f ./amp-agent && \
    go build -o /go/bin/amp-agent && \
    apk del build-deps && cd / && rm -rf $GOPATH/src /var/cache/apk/*

HEALTHCHECK --interval=5s --timeout=15s --retries=24 CMD curl localhost:3000/api/v1/health

CMD ["/go/bin/amp-agent"]

