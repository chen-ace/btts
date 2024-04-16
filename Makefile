# 生成Makefile，要生成windows平台、mac、linux平台的arm，x86架构的多个可执行文件
# 生成的可执行文件在bin目录下
export GO111MODULE=on
export GOPROXY=https://goproxy.io,direct
LDFLAGS := -s -w

GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)
GOARM := 7
GOMAIN := main.go
OUTDIR := bin
GOVERSION := $(shell go version | cut -d " " -f3 | tr -d 'go')

.PHONY: all windows linux linux-arm macos clean check-version check-go

all: check-go check-version windows linux linux-arm macos macos-arm64

check-go:
	@echo "Checking if Go is installed..."
	@if ! which go > /dev/null; then echo "Go is not installed. Please install it first."; exit 1; fi

check-version:
	@echo "Checking Go version..."
	@if [ $(GOVERSION) \< 1.22 ]; then echo "Go version is less than 1.22. Please upgrade."; exit 1; fi

windows:
	GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)"  -o $(OUTDIR)/btts-windows-amd64.exe

linux:
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)"  -o $(OUTDIR)/btts-linux-amd64

linux-arm:
	GOOS=linux GOARCH=arm GOARM=$(GOARM) go build -trimpath -ldflags "$(LDFLAGS)"  -o $(OUTDIR)/btts-linux-arm

macos:
	GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags "$(LDFLAGS)"  -o $(OUTDIR)/btts-macos-amd64

macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags "$(LDFLAGS)"  -o $(OUTDIR)/btts-macos-arm64

clean:
	rm -rf $(OUTDIR)/*