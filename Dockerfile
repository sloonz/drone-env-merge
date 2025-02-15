FROM golang:alpine as builder
RUN apk add --no-cache git ca-certificates
WORKDIR /go/pkg/drone-env-merge
COPY . .
RUN go build

FROM alpine:3
EXPOSE 80
ENV GODEBUG netdns=go
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/pkg/drone-env-merge/drone-env-merge /bin/
ENTRYPOINT ["/bin/drone-env-merge"]
