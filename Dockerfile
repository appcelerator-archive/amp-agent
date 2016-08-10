FROM appcelerator/alpine:3.4


RUN apk update
RUN apk add bash go bzr git mercurial subversion openssh-client ca-certificates && mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p /go/src/github.com/appcelerator/amp-agent /go/bin
WORKDIR /go/src/github.com/appcelerator/amp-agent
COPY ./ ./
RUN rm -rf ./vendor
RUN go get -u github.com/Masterminds/glide/...
RUN glide install
RUN go build -o /go/bin/amp-agent             

RUN chmod +x /go/bin/*

HEALTHCHECK --interval=3s --timeout=10s --retries=6 CMD curl localhost:3000/api/v1/health

CMD ["/go/bin/amp-agent"]

