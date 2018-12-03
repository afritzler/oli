FROM golang:1.11.2
WORKDIR /go/src/github.com/afritzler/oli
COPY . .
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/afritzler/oli/oli .
ENTRYPOINT ["/oli"]
