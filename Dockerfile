from golang:1.20-alpine as builder
copy . /src
workdir /src
run go build

from alpine:3.17
run apk --no-cache add ca-certificates
copy --from=builder /src/github-pr-creator /github-pr-creator
entrypoint ["/github-pr-creator"]
