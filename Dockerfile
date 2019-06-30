FROM golang:alpine
WORKDIR $GOPATH/src/duko/slimbin
COPY . $GOPATH/src/duko/slimbin
RUN go get . && CGO_ENABLED=0 go build -o main

FROM scratch
COPY --from=0 /go/src/duko/slimbin/main /main
EXPOSE 80
EXPOSE 8080
ENTRYPOINT ["/main"]