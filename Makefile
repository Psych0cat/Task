GO=go
GO111MODULE := on
export GO111MODULE
MOD := -mod=vendor
all: build
BUILD.go = $(GO) build -i
all: build
build:
	$(GO) build $(MOD) Router/main.go


