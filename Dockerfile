FROM alpine:3.7
ENTRYPOINT ["/bin/oli"]

COPY . /go/src/github.com/afritzler/oli
RUN apk --no-cache add -t build-deps build-base go git \
	&& apk --no-cache add ca-certificates \
	&& cd /go/src/github.com/afritzler/oli \
	&& export GOPATH=/go \
	&& go build -o /bin/oli \
	&& rm -rf /go \
	&& apk del --purge build-dep