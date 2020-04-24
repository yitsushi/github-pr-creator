from alpine:3.11 as builder
run apk add --no-cache go
copy src /src
workdir /src
run go build

from alpine:3.11
run apk --no-cache add ca-certificates
copy --from=builder /src/github-pr-creator /github-pr-creator
entrypoint ["/github-pr-creator"]
