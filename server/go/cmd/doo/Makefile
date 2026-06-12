# doo CLI 构建
VERSION ?= dev
PKG     := github.com/dootask/tools/server/go/cmd/doo/internal/commands
LDFLAGS := -s -w -X $(PKG).Version=$(VERSION)
DIST    := dist

.PHONY: build build-all test vet clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(DIST)/doo .

# 交叉编译矩阵：linux（容器目标）+ 开发机
build-all:
	CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/doo-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux  GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/doo-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/doo-darwin-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/doo-darwin-amd64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST)/doo-windows-amd64.exe .

test:
	go test ./...

vet:
	go vet ./...

clean:
	rm -rf $(DIST)
