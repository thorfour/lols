.PHONY: plugin server clean

GO_VERSION="1.11"
BIN="lols"

go = @docker run \
        --rm \
        -v ${PWD}:/go/src/github.com/thorfour/lols \
        -w /go/src/github.com/thorfour/lols \
        -u $$(id -u) \
        -e XDG_CACHE_HOME=/tmp/.cache \
        -e CGO_ENABLED=0 \
        -e GOOS=linux \
        -e GOPATH=/go \
        golang:$(GO_VERSION) \
        go

go_cgo = @docker run \
        --rm \
        -v ${PWD}:/go/src/github.com/thorfour/lols \
        -w /go/src/github.com/thorfour/lols \
        -u $$(id -u) \
        -e XDG_CACHE_HOME=/tmp/.cache \
        -e CGO_ENABLED=1 \
        -e GOOS=linux \
        -e GOPATH=/go \
        golang:$(GO_VERSION) \
        go


clean:
	rm -rf ./bin
	rm -f ca-certificates.crt
plugin: 
	mkdir -p ./bin/plugin
	$(go_cgo) build -buildmode=plugin -o ./bin/plugin/$(BIN) ./cmd/plugin/
server:
	mkdir -p ./bin/server
	$(go) build -o ./bin/server/$(BIN) ./cmd/server
docker:
	cp /etc/ssl/certs/ca-certificates.crt .
	docker build -t quay.io/thorfour/lols .
