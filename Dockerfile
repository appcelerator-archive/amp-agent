#FROM appcelerator/alpine:20160726
FROM golang:1.7-alpine

ENV GOPATH /go
ENV PATH $PATH:/go/bin

COPY ./ /go/src/github.com/appcelerator/amp-agent
RUN apk update && \
    apk --virtual build-deps add git make && \
    apk add curl bash && \

    cd /go/src/github.com/appcelerator/amp-agent && \
    go get -u github.com/Masterminds/glide/... && \
    glide install && \
    rm -f ./amp-agent && \
    make install && \
    rm /go/bin/glide && \
    apk del build-deps && cd / && rm -rf /go/src /go/pkg /var/cache/apk/* /root/.cache /root/.glide

HEALTHCHECK --interval=10s --timeout=15s --retries=12 CMD curl localhost:3000/api/v1/health

CMD ["/go/bin/amp-agent"]
