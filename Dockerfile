FROM --platform=$BUILDPLATFORM golang:alpine
WORKDIR $GOPATH/src/duko/slimbin
ARG TARGETOS
ARG TARGETARCH
COPY . $GOPATH/src/duko/slimbin
RUN --mount=type=cache,target=/root/.cache/go-build \
	--mount=type=cache,target=$GOPATH/pkg \
    go get . && GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -o main

FROM scratch
COPY --from=0 /go/src/duko/slimbin/main /main
EXPOSE 80
EXPOSE 8080
ENTRYPOINT ["/main"]
